import { Injectable } from '@nestjs/common';

import * as lark from '@larksuiteoapi/node-sdk';

@Injectable()
export class AppService {
  async getHello(): Promise<string> {
    const client = new lark.Client({
      appId: 'cli_a68feeb3967c900c',
      appSecret: 'VOKZ7eIM03sqqnhm7x9KHg5hh8eoTanc',
    });
    // 自建应用获取 tenant_access_token
    const res = await client.request({
      method: 'POST',
      url: 'https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal',
      data: {
        app_id: 'cli_a68feeb3967c900c',
        app_secret: 'VOKZ7eIM03sqqnhm7x9KHg5hh8eoTanc',
      },
    });

    if (res.code === 0) {
      // 成功获取token
      const tenantAccessToken = res.tenant_access_token;
      console.log('tenant_access_token:', tenantAccessToken);
      // 获取 user_access_token

      // 获取表格记录
      // {
      //   "code": 0,
      //   "data": {
      //     "sheets": [
      //       {
      //         "grid_properties": {
      //           "column_count": 20,
      //           "frozen_column_count": 0,
      //           "frozen_row_count": 0,
      //           "row_count": 200
      //         },
      //         "hidden": false,
      //         "index": 0,
      //         "resource_type": "sheet",
      //         "sheet_id": "3ce009",
      //         "title": "Sheet1"
      //       }
      //     ]
      //   },
      //   "msg": ""
      // }

      // const res11 = await client.request({
      //   method: 'GET',
      //   url: 'https://open.feishu.cn/open-apis/sheets/v3/spreadsheets/BeIRfbqHzl2crAdBv1gcIbYPnJO/sheets/:sheet_id',
      //   headers: {
      //     Authorization: 'Bearer ' + tenantAccessToken,
      //   },
      // });


      // 新建文件夹

      // token
      const res1 = await client.request({
        method: 'GET',
        url: 'https://open.feishu.cn/open-apis/drive/explorer/v2/root_folder/meta',
        headers: {
          Authorization: 'Bearer ' + tenantAccessToken,
        },
      });
      if (res1.code === 0) {
        console.error(res1.data);

        const ftoken = res1.data.token;

        const res3 = await client.request({
          method: 'POST',
          url: 'https://open.feishu.cn/open-apis/drive/v1/files/create_folder',
          headers: {
            Authorization: 'Bearer ' + tenantAccessToken,
          },
          data: {
            name: '配置项',
            folder_token: ftoken,
          },
        });
        if (res3.code === 0) {
          console.log(res3.data);
        } else {
          console.error('add 文件夹失败:', res3.msg);
        }
        // // 创建表格
        const res2 = await client.request({
          method: 'POST',
          url: 'https://open.feishu.cn/open-apis/sheets/v3/spreadsheets',
          data: {
            title: '测试',
            folder_token: ftoken,
          },
          headers: {
            Authorization: 'Bearer ' + tenantAccessToken,
          },
        });
        if (res2.code === 0) {
          console.error(res2.data);
        } else {
          console.error('add table:', res2.msg);
        }
      } else {
        console.error('add table:', res1.msg);
      }
    } else {
      console.error('Failed to get app access token:', res.msg);
    }
    console.error(res.code);
    return 'Hello World!';
  }
}
