import bcrypt
import streamlit as st
import streamlit_authenticator as stauth
from streamlit_authenticator.utilities.exceptions import LoginError


@st.cache_data(ttl="24h")
def get_credentials():
    cres = st.secrets.credentials
    # 定义字典，初始化字典
    users = {}
    for name, pwd in cres.to_dict().items():
        users[name] = {"name": name, 'password': bcrypt.hashpw(pwd["password"].encode(), bcrypt.gensalt()).decode()}
    return {'usernames': users}


def login():
    # 用户信息，后续可以来自DB
    credentials = get_credentials()
    authenticator = stauth.Authenticate(credentials, 'wechat-gptbot', 'qwertyuioasdfghjkl', cookie_expiry_days=30)
    try:
        authenticator.login()
    except LoginError as e:
        st.error(e)


def check_login():
    if st.session_state["authentication_status"]:
        # 已经登录了，校验公众号是否登录
        if st.session_state["wechat_user"]:
            return
    st.switch_page("主页.py")