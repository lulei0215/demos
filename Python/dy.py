import requests
import time
import json
import qrcode
import uuid
import base64


# 获取登录二维码 https://sso.douyin.com/get_qrcode
def get_qrcode():
    url = "https://sso.douyin.com/get_qrcode"
    headers = {}
    params = {
        "next": "https://creator.douyin.com/creator-micro/home",
        "aid": "2906",
        "service": "https://creator.douyin.com",
        "is_vcd": 1,
        "fp": "lrjypci4_wubS76S8_JwNQ_4AcP_9Hsm_xv3CEspCc06j",
        "msToken": "EPQa7rs032ANnXQnV5nmt3nHQ0lDsD-b8otCDH8jZk6ykVhBJjnxZ6GTruUFkMetr3QLCE-ILBXyT6WPSzwnbs0XcZ9SoiazHuO1t6XYLNkVbFMOAnpUKvfLlW4=",
        "X-Bogus": "DFSzswVuEBtANaDttiy6mfLNKBTM",
        "_signature": "_02B4Z6wo00001BGSUdQAAIDDQ21xzcMQ8wQRklVAAGHTLgbHtC8kUeIZK080FEYXIWbVk.GiGxNx3W5iXmutoU0o3XGIOyVv77eIP8YGz-EBT57Zyy2ciztISuDvA-F0CQqkiJo8LtnhEkKD83",
    }
    data = {}
    response = requests.get(
        url,
        headers=headers,
        data=data,
        params=params,
    )
    res = response.json()
    # print(res["data"]["qrcode"])
    img = f"data:image/png;base64,{res['data']['qrcode']}"
    if img:
        create_qc_code(img)
    else:
        print("二维码获取失败")
    # print(response.text)


# 获取评论列表
def get_comment_list():
    url = "https://creator.douyin.com/comment/list"
    headers = {}
    params = {
        "aid": "2906",
        "app_name": "aweme_creator_platform",
        "device_platform": "web",
        "referer": "https://www.google.com/",
        "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        "cookie_enabled": True,
        "screen_width": 1440,
        "screen_height": 900,
        "browser_language": "zh-CN",
        "browser_platform": "MacIntel",
        "browser_name": "Mozilla",
        "browser_version": "5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        "browser_online": "true",
        "timezone_name": "Asia/Shanghai",
        "item_id": "@ifRh7bNIQ0fxuMhG6Ljpv8p2pCyXbiD1K6olPaxtgJiEsRpc4ziTGnT3VKlkZZwk",
        "count": 10,
        "msToken": "XQlCnS9TgTxbAXJup1sgaUW2gOIhsK7HDcUmJzrZyrxIudJmhf9M8_SlXrOhBbkgPPuOI6w2xlG5v2_yRaWdIRXHJQUX4w6KdDMDWe2pxIPU-uuNlg==",
        "X-Bogus": "DFSzswVLtVzANy6JtimXbDLNKBT1",
        "_signature": "_02B4Z6wo000013aFncgAAIDAJHq90W8J7jN2hZlAALgZdOsNG0gLhwRLS6X1vQDmkG4dpdqtTOhyK4XdfCqL4hEvSkqw87XJIzaJQGvu-7wCOCnwZK05Mi2b-oBOkEH2SEti54tA.lJtOrjJ1f",
    }


# 7.生成二维码
def create_qc_code(url):
    # 将base64格式转url
    print(url)
    img = qrcode.make(url)
    # 可以使用img.save()保存图片到本地
    img.save("dy_login_code.png")


if __name__ == "__main__":
    get_qrcode()

# data:image/png;base64,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
# data:image/png;base64,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
