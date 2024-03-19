//微信小程序发货
public function sendwx($send_data = []){

// $out_trade_no = '',$item_desc= '',$logistics_type= '',$tracking_no= '',$express_company= '',$item_desc= '',$consignor_contact= ''
// $out_trade_no,$item_desc,$logistics_type,$tracking_no,$express_company,$item_desc,$consignor_contact


$plateform_config_model = new WeappConfigModel();
$plateform_config = $plateform_config_model->getWeappConfig(1);
$config = [
'app_id' => $plateform_config['data']['value']["appid"] ?? '',
'secret' => $plateform_config['data']['value']["appsecret"] ?? '',
'token' => $plateform_config['data']['value']["token"] ?? '',
'aes_key' => $plateform_config['data']['value']["config_key"] ?? '',
'log' => [
'level' => 'debug',
'permission' => 0777,
'file' => 'runtime/log/wechat/oplatform.logs',
],
];

$app =Factory::miniProgram($config);

$accessToken = $app->access_token; // EasyWeChat\Core\AccessToken 实例
$token = $accessToken->getToken(); // token 字符串
$token = $accessToken->getToken(true);

// $post_url = 'https://api.weixin.qq.com/product/delivery/send?access_token='.$token['access_token'];
$post_url = 'https://api.weixin.qq.com/wxa/sec/order/upload_shipping_info?access_token='.$token['access_token'];

$post_data = array(
"order_key"=>[
"order_number_type"=>2, //微信支付单号
"out_trade_no"=>$send_data['out_trade_no'], //原支付交易对应的微信订单号
"mchid"=>"1659577072"
],
"delivery_mode"=>1,
"logistics_type"=>$send_data['logistics_type'],//1、实体物流配送采用快递公司进行实体物流配送形式 2、同城配送 3、虚拟商品，虚拟商品，例如话费充值，点卡等，无实体配送形式 4、用户自提
"shipping_list"=>[
[
"tracking_no"=>$send_data['tracking_no'], //物流单号
"express_company"=>$send_data['express_company'],//物流公司编码
"item_desc"=>$send_data['item_desc'],//商品信息
"contact"=>[
"consignor_contact"=>$send_data['consignor_contact'] //寄件人联系方式
]
]
],
"upload_time"=>date(DATE_RFC3339),
"payer"=>[
"openid"=>$send_data['openid']//用户标识
]
);
$post_data['logistics_type'] = $send_data['logistics_type'];

if($post_data['logistics_type'] == 1){
$post_data['shipping_list'][0] =[
"tracking_no"=>$send_data['tracking_no'], //物流单号
"express_company"=>$send_data['express_company'],//物流公司编码
"item_desc"=>$send_data['item_desc'],//商品信息
"contact"=>[
"consignor_contact"=>$send_data['consignor_contact'] //寄件人联系方式
]
] ;
}

$client = new \GuzzleHttp\Client();
$response = $client->request('POST', $post_url, [
'json' => $post_data
]);

$result = $response->getBody()->getContents();
return $result;
}