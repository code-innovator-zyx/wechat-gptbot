import streamlit as st
from apis import *
res = check_login()
if res["code"] == 511:
    st.switch_page("home.py")
st.set_page_config(page_title='提示词管理', page_icon='😊', layout='wide',
                   initial_sidebar_state="expanded")

st.title('🐂 聊天提示词管理')
