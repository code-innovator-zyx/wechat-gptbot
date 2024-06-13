import streamlit as st

st.set_page_config(page_title='插件管理', page_icon='🔩', layout='wide',
                   initial_sidebar_state="expanded")

st.title('🚗 机器人儿插件管理')
step, weather_forecast, hot_search = st.tabs(["微信运动", "天气预报", "每日热搜"])
with step:
    st.info("TODO")

with weather_forecast:
    st.info("TODO")

with hot_search:
    st.info("TODO")
