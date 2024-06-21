import time

from st_pages import show_pages_from_config, hide_pages

show_pages_from_config()

import streamlit as st

from apis import *

st.set_page_config(page_title='wechat-gptbot', page_icon='🤖', layout='wide',
                   initial_sidebar_state="expanded")
st.warning("️当前页面的一切修改都是临时的，服务重启会重置", icon='⚠️')

# 创建三个列 为了 让其居中对齐
col1, col2, col3 = st.columns(3)
with col2:
    st.markdown(
        f"<h1 style='text-align: center; color: green; font-size:28px;'>微信gpt机器人🤖</h1>",
        unsafe_allow_html=True)
    res = check_login()
    if res["code"] == 511:
        st.session_state["login"] = "failed"
        hide_pages(["主页", "模型配置", "提示词配置"])
        st.image(res["data"]["qr_url"], caption="扫描二维码登录微信机器人")
        time.sleep(2)  # 2秒刷新一次
        st.rerun()
    else:
        st.session_state["login"] = "success"
        # 展示用户名
        st.markdown(
            f"<h1 style='text-align: center; color: green; font-size:28px;'>【{res['data']['user_name']}】登录成功,现在您"
            f"可以对配置进行动态调整啦 </h1>",
            unsafe_allow_html=True)
