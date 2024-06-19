import pandas as pd
import streamlit as st
from util import *

res = check_login()
if res["code"] == 511:
    st.switch_page("home.py")
st.set_page_config(page_title='插件管理', page_icon='🔩', layout='wide',
                   initial_sidebar_state="expanded")

st.title('🚗 机器人儿插件管理')
step, weather_forecast, hot_search = st.tabs(["微信运动", "天气预报", "每日热搜"])
with step:
    st.info("TODO")


def comment_keyword_onlisten():
    editors = st.session_state.get("friends_data_editor", None)
    if editors:
        for i, editor in editors.get("edited_rows", {}).items():
            if len(editor) != 0:
                delete_weather_receiver(current_set['users'][i]["name"])
                st.success(f"delete {current_set['users'][i]}", icon="✅")


with weather_forecast:
    st.info("准点天气推送，每天早上八点推送天气给指定用户")
    frineds = get_friends()
    current_set = get_weather_cron_setting()
    df = pd.DataFrame(current_set["users"])
    df["check"] = True
    st.data_editor(df, key="friends_data_editor", column_order=[
        "name",
        "city",
        "check"
    ], column_config={
        "name": st.column_config.TextColumn("微信昵称(或备注)"),
        "city": st.column_config.TextColumn("关联查询城市"),
        "check": st.column_config.CheckboxColumn("已绑定"),
    }, on_change=comment_keyword_onlisten, hide_index=True, use_container_width=True)

    with st.expander("添加接收人或修改城市"):
        with st.form("接收人"):
            receiver = st.selectbox("选择接收人", frineds["data"]["users"])
            city = st.text_input("请准确填写城市，如'成都'")
            if st.form_submit_button("submit"):
                add_weather_receiver(receiver, city)
                st.info("成功")
with hot_search:
    st.info("TODO")
