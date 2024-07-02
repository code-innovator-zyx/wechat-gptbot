import time

import streamlit as st

import auth
from apis import *

st.set_page_config(page_title='wechat-gptbot', page_icon='🤖', layout='wide',
                   initial_sidebar_state="expanded")

auth.login()
# 创建三个列 为了 让其居中对齐
if st.session_state["authentication_status"]:
    col1, col2, col3 = st.columns(3)
    with col2:
        st.markdown(
            f"<h1 style='text-align: center; color: green; font-size:28px;'>微信gpt机器人🤖</h1>",
            unsafe_allow_html=True)
        res = check_login()
        if res["code"] == 511:
            st.session_state["wechat_user"] = False
            st.image(res["data"]["qr_url"], caption="扫描二维码登录微信机器人")
            time.sleep(2)  # 2秒刷新一次
            st.rerun()
        else:
            st.session_state["wechat_user"] = True
            # 展示用户名
            st.markdown(
                f"<h1 style='text-align: center; color: green; font-size:28px;'>【{res['data']['user_name']}】登录成功,现在您"
                f"可以对配置进行动态调整啦 </h1>",
                unsafe_allow_html=True)
elif st.session_state["authentication_status"] is False:
    st.error('账号密码错误')
elif st.session_state["authentication_status"] is None:
    st.warning('请输入账号密码登录')
