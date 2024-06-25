import os
import secrets
from typing import Dict, Any

import requests
import streamlit
from streamlit import secrets

api_port = os.getenv("APIPORT", 8502)


# 调用本地api


def get_routers(name, group="router"):
    return secrets.get(group, {}).get(name, None)


def check_login():
    response = call_api("get", "check_login")
    return response.json()


def get_models():
    response = call_api("get", "current_models")
    return response.json()


@streamlit.cache_data(ttl="1h")
def get_friends():
    response = call_api("get", "friends")
    return response.json()


def reset_models(text_model):
    response = call_api("post", "reset_model", json={"text_model": text_model})
    return response.status_code


def get_cron_setting(plugin_name):
    response = call_api("get", f"{plugin_name}_cron_setting")
    return response.json()


def delete_receiver(plugin_name, name):
    call_api("delete", f"{plugin_name}_receiver", params={"name": name})


def reset_receiver(plugin_name, args):
    response = call_api("post", f"{plugin_name}_receiver", json=args)
    return response.json()


def text_models():
    return ["gpt-3.5-turbo", "gpt-4", "gpt-4-0125-preview", "gpt-4-turbo", "gpt-4o"]


def drawing_models():
    return ["dall-e-2", "dall-e-3"]


def reset_cron(plugin_name, desc):
    response = call_api("post", "reset_cron", json={"plugin_name": plugin_name, "desc": desc
                                                    })
    return response.json()


def reset_rss(source, top_n):
    response = call_api("post", "reset_rss", json={"source": source, "top_n": top_n})
    return response.json()


def call_api(method, router_name, params: Dict[str, Any] = None,
             json: Dict[str, Any] = None, ):
    router = get_routers(router_name)
    if router is None:
        return
    url = f'http://127.0.0.1:{api_port}{router}'
    kws = {}
    if params:
        kws["params"] = params
    if json:
        kws["json"] = json
    # Using the requests library for synchronous HTTP requests
    response = requests.request(method, url, **kws)
    return response
