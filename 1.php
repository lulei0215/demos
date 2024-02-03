<?php

namespace app\adminapi\controller;

use app\admin\model\Admin;
use app\admin\model\Admin as AdminModle;
use app\admin\model\AdminLog;
use app\admin\model\Buying;
use app\admin\model\mall\OrderProduct;
use app\admin\model\mall\Product;
use app\admin\model\mall\Skuvalue;
use app\common\controller\Backend as Api;
use app\common\library\Ems;
use app\common\library\Sms;
use app\admin\model\mall\Bill;
use app\admin\model\mall\Order as MallOrder;
use app\common\model\shop\order\After;
use fast\Random;
use think\Config;
use think\Db;
use think\Exception;
use think\exception\PDOException;
use think\exception\ValidateException;
use think\Hook;
use think\Request;
use think\Session;
use think\Validate;
use app\api\controller\Exports;

/**
 * 订单
 */
class Order extends Api
{
    protected $noNeedLogin = ['*'];
    protected $noNeedRight = ['*'];

    public function _initialize()
    {
        parent::_initialize();
        //        $this->Adminmodel = model('Admin');
        if (!Config::get('fastadmin.usercenter')) {
            $this->error(__('User center already closed'));
        }
        $this->model = new \app\admin\model\mall\Order;
    }

    function combineStrings($str1, $str2)
    {
        $length1 = strlen(preg_replace("/[\x{4e00}-\x{9fa5}]/u", "**", $str1));
        $length2 = strlen(preg_replace("/[\x{4e00}-\x{9fa5}]/u", "**", $str2));

        if ($length1 + $length2 > 32) {
            return  $str1 . "<BR><RIGHT>" . $str2 . "</RIGHT><BR>";
        }

        $spacesNeeded = 32 - $length1 - $length2;
        $spaces = str_repeat(' ', $spacesNeeded);

        return $str1 . $spaces . $str2;
    }
    /**
     * 查看
     */
    public function list()
    {

        $shop_id = gettoken($this->request->header('token'));
        if ($shop_id['code'] != 1) {
            return json($shop_id);
        }
        $shop_id = $shop_id['data']['shop_id'];

        $limit = $this->request->post("limit", 10);
        $ing = $this->model
            ->where("shop_id", $shop_id)
            ->where("order_id", 0)
            ->where("print_status", 1)
            ->count();
        $finish = $this->model
            ->where("shop_id", $shop_id)
            ->where("order_id", 0)
            ->where("print_status", 2)
            ->count();
        $where = [];
        $where['print_status'] = ['>', "0"];
        if ($this->request->param("desk_id", '')) {
            $where['desk_id'] = $this->request->param('desk_id');
        }
        if ($this->request->param("start_time", '') && $this->request->param("end_time", '')) {
            $where['createtime'] = ['>', strtotime($this->request->param("start_time"))];
            $where['createtime'] = ['<', strtotime($this->request->param("end_time"))];
        }
        if ($this->request->param("status", '')) {
            $where['print_status'] = $this->request->param('status');
        }



        $orderModel = new \app\admin\model\mall\Order();
        $result = $orderModel->with(["orderproduct.product", "orderproduct.sku", 'desk', 'shop'])
            ->where("shop_id", $shop_id)
            ->where("order_id", 0)
            ->order("id desc")
            ->paginate($limit);




        foreach ($result as $key => $value) {

            $num = 0;
            $p_total = $value['total_price'];

            $add_orders = $orderModel->with(["orderproduct.product", "orderproduct.sku", 'desk', 'shop'])
                ->where("order_id", $value['id'])
                ->order("id desc")
                ->select();

            foreach ($value['orderproduct'] as $pkey => $pvalue) {
                if ($pvalue['sku'] != NULL) {
                    $num += $pvalue['number'];
                } else {
                    $num += $pvalue['number'];
                }
            }

            if (count($add_orders) > 0) {


                $result[$key]['add_orders'] = $add_orders;

                foreach ($add_orders as $pkey1 => $pvalue1) {

                    foreach ($pvalue1['orderproduct'] as $pkey11 => $pvalue11) {
                        if ($pvalue11['sku'] != NULL) {
                            $num += $pvalue11['number'];
                        } else {
                            $num += $pvalue11['number'];
                        }
                    }



                    $p_total += $pvalue1['total_price'];
                }
            } else {

                $result[$key]['add_orders'] = [];
            }


            $result[$key]['total'] = $p_total;
            $result[$key]['num'] = $num;


            $result[$key]['id'] = (string)$result[$key]['id'];

            $value['id'] = (string)$value['id']; // 整形数字太大js会失准
            $value['have_paid_status'] = $value['have_paid'];
            $value['have_received_status'] = $value['have_received'];
            $value['have_commented_status'] = $value['have_commented'];
            $value['have_made_status'] = $value['have_made'];
        }


        $result = [
            'code' => 1,
            'msg' => '成功',
            'data' => ["list" => $result, 'ing' => $ing, "finish" => $finish],
        ];
        return json($result);
    }
    /**
     * 查看
     */
    public function detail()
    {

        $shop_id = gettoken($this->request->header('token'));
        if ($shop_id['code'] != 1) {
            return json($shop_id);
        }
        $shop_id = $shop_id['data']['shop_id'];

        $type = $this->request->post("type", 1);
        $id = $this->request->post("id");

        if ($type == 1) {

            $status = $orderModel
                ->with([
                    "orderproduct.sku", 'desk',
                    'user' => function ($query) {
                        $query->field('id,username');
                    }
                ])
                ->find($order_id);
            $num = 0;
            foreach ($status['orderproduct'] as $pkey1 => $pvalue1) {
                $num += $pvalue1['number'];
            }

            $p_total =  $status['total_price'];

            $status["add_orders"] = [];
            if ($status['order_id'] == 0) {
                $orders = $orderModel
                    ->with(["orderproduct.sku", 'desk'])
                    ->where("order_id", $status['id'])
                    ->order('createtime desc')
                    ->select();
                if (count($orders) > 0) {


                    foreach ($orders as $key => $value) {
                        foreach ($value['orderproduct'] as $pkey => $pvalue) {
                            $num += $pvalue['number'];
                        }

                        $p_total +=  $value['total_price'];
                    }

                    $status["add_orders"] = $orders;
                }
            }
            $status['num'] = $num;
            $status['total'] = $p_total;

            if (!$status) {
                $result = [
                    'code' => 0,
                    'msg' => '订单未找到',
                    'data' => ''
                ];
                return json($result);
            }
        } else {
            $info =    Buying::with(["desk", "gettype"])->where("id", $id)->find();
            if (!$info) {
                $info = [];
            } else {
                $info['createtime'] = date("Y-m-d H:i:s", $info['createtime']);
                $info['updatetime'] = date("Y-m-d H:i:s", $info['updatetime']);
            }
        }



        $result = [
            'code' => 1,
            'msg' => '成功',
            'time' => Request::instance()->server('REQUEST_TIME'),
            'data' => $info,
        ];
        return json($result);
    }
    public function add()
    {
        $shop_id = gettoken($this->request->header('token'));
        if ($shop_id['code'] != 1) {
            return json($shop_id);
        }
        $shop_id = $shop_id['data']['shop_id'];
    }


    /**
     * @return 退菜
     */
    public function back()
    {
        //        $id = $this->request->post('id');
        $shop_id = gettoken($this->request->header('token'));
        if ($shop_id['code'] != 1) {
            return json($shop_id);
        }
        $shop_id = $shop_id['data']['shop_id'];
        $product_id = $this->request->post('ids');
        $number = $this->request->post('number');
        $order_id = $this->request->post('order_id');



        $productIds = explode(',', $product_id);
        $numbers = explode(',', $number);

        $products = [];
        $orderProduct = [];
        $order_product = [];
        $price = 0;
        //        退菜    菜id 菜数量 费用
        foreach ($productIds as $key => &$productId) {

            $value = OrderProduct::find($productId);
            if (!$value) {
                return json_encode([
                    'code' => 0,
                    'msg' => '商品信息未查到',
                    'data' => []
                ]);
            }

            //退掉 所有数量
            if ($value['number'] == $numbers[$key]) {

                $price += $numbers[$key] * $value['price'];
                $ss =  Db::name("mall_order_product")->where('id', $productId)->update(['status' => -1]);
            } else {

                if ($value['number'] < $numbers[$key]) {
                    return json_encode([
                        'code' => 0,
                        'msg' => '退菜数量超过了 订单数量',
                        'data' => []
                    ]);
                    die();
                }
                //退掉 部分数量
                OrderProduct::where('id', $productId)->setDec('number', $numbers[$key]);
                $order_product[$key] = [
                    'order_id' => $value['order_id'],
                    'user_id' => $value['user_id'],
                    'product_id' => $value['product_id'],
                    'number' => $numbers[$key],
                    'title' => $value['title'],
                    'tip' => $value['tip'],
                    'image' => $value['image'],
                    'spec' => $value['spec'],
                    'price' => $value['price'],
                    // 'total_price' => $value['price'] * $numbers[$key],
                    'status' => -1,
                    'createtime' => time(),
                    'updatetime' => time(),
                ];
                $price += $numbers[$key] * $value['price'];
            }
        }

        if (count($order_product) > 0) {
            (new OrderProduct)->insertAll($order_product);
        }

        if ($price > 0) {
            (new \app\admin\model\mall\Order())->where('id', $order_id)->setDec('order_price', $price);
            (new \app\admin\model\mall\Order())->where('id', $order_id)->setDec('total_price', $price);
        }
        $result = [
            'code' => 1,
            'msg' => '退菜成功',
            'time' => Request::instance()->server('REQUEST_TIME'),
            'data' => [],
        ];
        return json($result);
    }
    //    待处理订单
    public function wait()
    {

        $shop_id = gettoken($this->request->header('token'));
        if ($shop_id['code'] != 1) {
            return json($shop_id);
        }
        $shop_id = $shop_id['data']['shop_id'];
        $list = $this->model
            ->with([
                //                'shop',
                'orderproduct',
                'desk',
                //                'user' => function ($query) {
                //                    $query->field('id,username');
                //                }
            ])
            ->where("status", 1)->where("print_status", 0)
            ->where("shop_id", $shop_id)
            ->order("id desc")
            ->select();
        if (count($list) > 0) {
            foreach ($list as &$item) {
                $item['id'] = (string)$item['id']; // 整形数字太大js会失准
                $item['have_paid_status'] = $item['have_paid'];
                $item['have_received_status'] = $item['have_received'];
                $item['have_commented_status'] = $item['have_commented'];
                $item['have_made_status'] = $item['have_made'];
                $item['order_type'] = 1;
            }
        }

        $buyorder = Buying::with(["gettype", 'shop', 'desk'])
            ->where("shop_id", $shop_id)
            ->where("status", 0)
            ->order("id desc")
            ->select();
        if (count($buyorder) > 0) {
            foreach ($buyorder as  $item) {
                $item['order_type'] = 2; // 整形数字太大js会失准
                // $item['createtime'] = date("Y-m-d H:i:s",$item['createtime']);
            }
        }
        $ss = array_merge($list, $buyorder);
        $result = [
            'code' => 1,
            'msg' => '成功',
            'data' => $ss,
        ];
        return json($result);
    }
    //    打印小票
    public function print()
    {
        $shop_id = gettoken($this->request->header('token'));
        if ($shop_id['code'] != 1) {
            return json($shop_id);
        }
        $shop_id = $shop_id['data']['shop_id'];
        $order_id = $this->request->post("order_id", "");
        if (!$order_id) {
            return json_encode([
                'code' => 0,
                'msg' => '订单号有误',
                'data' => []
            ]);
        }
        $total = $this->model->find($order_id);
        if (!$total) {
            return json_encode([
                'code' => 0,
                'msg' => '订单号未查到',
                'data' => []
            ]);
        }

        $this->print_s($order_id);
        $result = [
            'code' => 1,
            'msg' => '成功',
            'data' => "",
        ];
        return json($result);
    }
    //    上菜
    public function onfood()
    {
        $shop_id = gettoken($this->request->header('token'));
        if ($shop_id['code'] != 1) {
            return json($shop_id);
        }
        $shop_id = $shop_id['data']['shop_id'];
        $order_id = $this->request->post("order_id", "");
        if (!$order_id) {
            return json_encode([
                'code' => 0,
                'msg' => '订单号有误',
                'data' => []
            ]);
        }
        $total = $this->model->find($order_id);
        if (!$total) {
            return json_encode([
                'code' => 0,
                'msg' => '订单号未查到',
                'data' => []
            ]);
        }
        $total->print_status = 1;
        $total->updatetime = time();
        $total->save();
        $result = [
            'code' => 1,
            'msg' => '成功',
            'data' => "",
        ];
        return json($result);
    }
    public function print_s($order_id)
    {

        $orderModel = new \app\admin\model\mall\Order();
        $order = $orderModel->with(["orderproduct.sku.skuvalue", 'desk', 'shop'])->where('id', $order_id)->find();
        if (!$order) {
            return false;
        }
        $shop_id = $order['shop_id'];
        $print = new Exports();
        if ($order['order_id'] == 0) {
            $text =      "<C><font# bolder=1 height=2 width=2>新订单</font#></C><BR>";
        } else {
            $text =      "<C><font# bolder=1 height=2 width=2>加菜订单</font#></C><BR>";
        }


        $text .=  "<C><font# bolder=1 height=2 width=2>桌号：" . $order['desk']['title'] . "</font#></C><BR>";

        $text .= "<LEFT><font# bolder=0 height=2 width=1>备注: " . $order["remark"] . "</font#></LEFT><BR>";

        $text .= "<C>********************************</C><BR><LEFT>订单编号: " . $order["out_trade_no"] . "</LEFT><BR><LEFT>下单时间: " . $order['createtime'] . "</LEFT><BR><C>--------------商品--------------</C><BR>";


        foreach ($order['orderproduct'] as $k => $v) {
            if ($v['sku']) {
                $text .=      $this->combineStrings($v['title'] . $v['sku']['value'], 'x' . $v['number']);
            } else {
                $text .=      $this->combineStrings($v['title'], 'x' . $v['number']);
            }
        }
        $text .=  "<C>--------------------------------</C><BR>";
        $sns = Db::name("mall_printers")->where("shop_id", $shop_id)->select();
        if (count($sns) == 0) {
            $result = [
                'code' => 0,
                'msg' => '没有打印机，请添加',
                'data' => "",
            ];
            return json($result);
        }

        for ($i = 0; $i < count($sns); $i++) {
            $ss = $print->print($sns[$i]['sn'], $text);
            if ($ss['code'] == 0) {
                return json($ss);;
            }
        }

        $result = [
            'code' => 1,
            'msg' => '打印结束',
            'data' => "",
        ];
        return json($result);
    }




    //结束订单
    public function finish()
    {
        $shop_id = gettoken($this->request->header('token'));
        if ($shop_id['code'] != 1) {
            return json($shop_id);
        }
        $shop_id = $shop_id['data']['shop_id'];
        $order_id = $this->request->post("order_id", "");
        if (!$order_id) {
            return json_encode([
                'code' => 0,
                'msg' => '订单号有误',
                'data' => []
            ]);
        }
        $total = $this->model->find($order_id);
        if (!$total) {
            return json_encode([
                'code' => 0,
                'msg' => '订单号未查到',
                'data' => []
            ]);
        }
        $total->print_status = 2;
        $total->endtime = time();
        $total->save();


        $result = [
            'code' => 1,
            'msg' => '成功',
            'time' => Request::instance()->server('REQUEST_TIME'),
            'data' => "",
        ];
        return json($result);
    }
    //        今日订单数据
    public function today()
    {
        $shop_id = gettoken($this->request->header('token'));
        if ($shop_id['code'] != 1) {
            return json($shop_id);
        }
        $shop_id = $shop_id['data']['shop_id'];
        $todayCount = $this->model
            ->where("shop_id", $shop_id)
            ->where("status", 1)
            ->whereTime('createtime', 'today')
            ->count();
        $yesterdayCount = $this->model
            ->where("shop_id", $shop_id)
            ->where("status", 1)
            ->whereTime('createtime', 'yesterday')
            ->count();

        $t_todayCount = Buying::whereTime('createtime', 'today')->where("shop_id", $shop_id)->count();
        $t_yesterdayCount = Buying::whereTime('createtime', 'yesterday')->where("shop_id", $shop_id)->count();


        $result = [
            'code' => 1,
            'msg' => '成功',
            'time' => Request::instance()->server('REQUEST_TIME'),
            'data' => [
                "tuangou" => [
                    'todayCount' => $todayCount,
                    'yesterdayCount' => $yesterdayCount
                ],
                "diancan" => [
                    "todayCount" => $t_todayCount,
                    "yesterdayCount" => $t_yesterdayCount
                ]
            ]

        ];
        return json($result);
    }


    //       店铺数据 今日 本周 本月
    // 消费用户 营业总额 卡券订单  团购 点餐
    public function datas()
    {
        $shop_id = gettoken($this->request->header('token'));
        if ($shop_id['code'] != 1) {
            return json($shop_id);
        }
        $shop_id = $shop_id['data']['shop_id'];

        $type = $this->request->post('type', 'today');

        $bill = new Bill();
        $order  = new MallOrder();
        $buying = new Buying();

        $yesterdayCount = 0;

        $todayCount = $bill->where("shop_id", $shop_id)
            ->where("type", 1)
            ->whereTime('createtime', 'today')
            ->sum("total_price");

        switch ($type) {
            case  "today":
                $ren = $order->where("shop_id", $shop_id)
                    ->where("type", 1)
                    ->whereTime('createtime', 'today')
                    ->group("user_id")
                    ->count();

                $billCount = $order->where("shop_id", $shop_id)
                    ->where("type", 1)
                    ->whereTime('createtime', 'today')
                    ->sum("total_price");

                $yesterdayCount = $order->where("shop_id", $shop_id)
                    ->where("type", 1)
                    ->whereTime('createtime', 'yesterday')
                    ->sum("total_price");

                $orderCount = $order->where("shop_id", $shop_id)
                    ->where("status", 1)
                    ->whereTime('createtime', 'today')
                    ->count();
                $buyingCount = $buying->where("shop_id", $shop_id)
                    ->where("status", 1)
                    ->whereTime('createtime', 'today')
                    ->count();

                break;
            case    "week":
                $ren = $bill->where("shop_id", $shop_id)
                    ->where("type", 1)
                    ->whereTime('createtime', 'week')
                    ->group("user_id")
                    ->count();
                $billCount = $bill->where("shop_id", $shop_id)
                    ->where("type", 1)
                    ->whereTime('createtime', 'week')
                    ->sum("total_price");
                $orderCount = $order->where("shop_id", $shop_id)
                    ->where("status", 1)
                    ->whereTime('createtime', 'week')
                    ->count();
                $buyingCount = $buying->where("shop_id", $shop_id)
                    ->where("status", 1)
                    ->whereTime('createtime', 'week')
                    ->count();
                break;

            default:
                $ren = $order->where("shop_id", $shop_id)
                    ->where("type", 1)
                    ->whereTime('createtime', 'month')
                    ->group("user_id")
                    ->count();
                $billCount = $order->where("shop_id", $shop_id)
                    ->where("type", 1)
                    ->whereTime('createtime', 'month')
                    ->sum("total_price");
                $orderCount = $order->where("shop_id", $shop_id)
                    ->where("status", 1)
                    ->whereTime('createtime', 'month')
                    ->count();
                $buyingCount = $buying->where("shop_id", $shop_id)
                    ->where("status", 1)
                    ->whereTime('createtime', 'month')
                    ->count();
        }

        $result = [
            'code' => 1,
            'msg' => '成功',
            'time' => Request::instance()->server('REQUEST_TIME'),
            'data' => [
                "today:今日营业额，yesterday昨日；ren消费用户，billcount营业总额，orderCount点餐订单，buyingCount团购订单",
                "today" => $todayCount,
                'yesterday' => $yesterdayCount,
                'ren' => $ren,
                'billcount' => $billCount,
                'orderCount' => $orderCount,
                "buyingCount" => $buyingCount
            ]

        ];
        return json($result);
    }
    public function charts()
    {
        $shop_id = gettoken($this->request->header('token'));
        if ($shop_id['code'] != 1) {
            return json($shop_id);
        }

        $shop_id = $shop_id['data']['shop_id'];
        $type = $this->request->post('type', 'day');
        $order  = new MallOrder();
        $data = [];
        $dates = [];
        $j = 0;
        switch ($type) {
            case  "day":

                for ($i = 6; $i >= 0; $i--) {

                    $date = date('Y-m-d', strtotime("-$i days"));
                    $data[$j] = $order->where("shop_id", $shop_id)
                        ->where("type", 1)
                        ->where('createtime', 'between', [strtotime($date . ' 00:00:00'), strtotime($date . ' 23:59:59')])
                        ->sum("total_price");
                    $dates[$j] = $date;
                    $j++;
                }

                break;
            case    "week":
                for ($i = 29; $i >= 0; $i--) {
                    $date = date('Y-m-d', strtotime("-$i days"));
                    $data[$j] = $order->where("shop_id", $shop_id)
                        ->where("type", 1)
                        ->where('createtime', 'between', [strtotime($date . ' 00:00:00'), strtotime($date . ' 23:59:59')])
                        ->sum("total_price");
                    $dates[$j] = $date;
                    $j++;
                }
                break;
            default:
                for ($i = 0; $i < 12; $i++) {
                    $date = date('Y-m', strtotime("+$i months"));
                    $data[$j] = $order->where("shop_id", $shop_id)
                        ->where("type", 1)
                        ->where('createtime', 'between', [strtotime($date . '-1 00:00:00'), strtotime(date('Y-m-t', strtotime($date)))])
                        ->sum("total_price");
                    $dates[$j] = substr($date, 5, 2);
                    $j++;
                }
        }

        $result = [
            'code' => 1,
            'msg' => '成功',
            'time' => Request::instance()->server('REQUEST_TIME'),
            'data' => [
                "data" => $data,
                "dates" => $dates
            ]

        ];
        return json($result);
    }
}
