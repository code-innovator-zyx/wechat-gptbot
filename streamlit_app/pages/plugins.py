import pandas as pd
import streamlit as st
from apis import *

res = check_login()
if res["code"] == 511:
    st.switch_page("home.py")
st.set_page_config(page_title='定时任务管理', page_icon='🔩', layout='wide',
                   initial_sidebar_state="expanded")

st.title('🚗 机器人儿定时任务管理')
step, weather_forecast, hot_search = st.tabs(["微信运动", "天气预报", "媒体新闻推送"])
frineds = get_friends()

weather_plugin_name = "WeatherPlugin"
news_plugin_name = "NewsPlugin"
sport_plugin_name = "StepPlugin"


@st.experimental_dialog("新增账号绑定")
def update_sport_receiver(current_users):
    receiver = st.selectbox("选择接收人", list(set(frineds["data"]["users"]) - set(current_users)))
    account = st.text_input("zepplife 注册的账号")
    pwd = st.text_input("zepplife 注册的账号密码", type="password")
    min_step = st.number_input("微信运动最小步数", min_value=10000, max_value=30000, step=5000)
    max_step = st.number_input("微信运动最大步数", min_value=30000, max_value=80000, step=5000)
    if st.button("确认修改", use_container_width=True, type="primary"):
        if reset_receiver(plugin_name="sport", args={"name": receiver, "account": account,
                                                     "pwd": pwd, "min": min_step, "max": max_step}).get("msg",
                                                                                                        "") == "ok":
            st.rerun()


@st.experimental_dialog("修改执行时间")
def update_cron_spec(plugin_name, current):
    desc = st.text_input("时间描述：", value=current, max_chars=50)
    if st.button("重置", type="primary", use_container_width=True, disabled=current == desc):
        if reset_cron(plugin_name, desc).get("msg",
                                             "") == "ok":
            st.rerun()
        else:
            st.warning("修改失败")


def sport_forecast_onlisten():
    editors = st.session_state.get("sport_data_editor", None)
    if editors:
        for i, editor in editors.get("edited_rows", {}).items():
            if len(editor) != 0:
                delete_receiver("sport", sport_set['users'][i]["name"])
                st.success(f"delete {sport_set['users'][i]}", icon="✅")


with step:
    st.info("微信装逼神器，懂的都懂")
    st.page_link("https://github.com/code-innovator-zyx/wechat-gptbot/blob/main/core/plugins/wechatMovement/README.md",
                 label="请先点击了解详情", icon="🏃🏻‍♀️")

    sport_set = get_cron_setting("sport")
    df = pd.DataFrame()
    if sport_set["users"]:
        if st.button("微信运动：" + sport_set["cron"]):
            update_cron_spec(sport_plugin_name, sport_set["cron"])
        df = pd.DataFrame(sport_set["users"])
        df["check"] = True
        # 将密码字段显示为隐藏字符
        df["pwd"] = df["pwd"].apply(lambda x: "*" * len(x))
        st.data_editor(df, key="sport_data_editor", column_order=[
            "name", "account", "pwd", "min", "max", "check"
        ], column_config={
            "name": st.column_config.TextColumn("微信昵称(或备注)"),
            "account": st.column_config.TextColumn("zepplife 账号"),
            "pwd": st.column_config.TextColumn("zepplife 密码"),
            "max": st.column_config.NumberColumn("设置微信最小步数"),
            "min": st.column_config.NumberColumn("设置微信最大步数"),
            "check": st.column_config.CheckboxColumn("是否绑定"),
        }, on_change=sport_forecast_onlisten, hide_index=True, use_container_width=True)
    if st.button("绑定账号", type="primary", use_container_width=True):
        update_sport_receiver(df["name"].tolist() if len(df) != 0 else [])


def weather_forecast_onlisten():
    editors = st.session_state.get("friends_data_editor", None)
    if editors:
        for i, editor in editors.get("edited_rows", {}).items():
            if len(editor) != 0:
                delete_receiver("weather", weather_set['users'][i]["name"])
                st.success(f"delete {weather_set['users'][i]}", icon="✅")


@st.experimental_dialog("新增用户")
def update_weather_receiver(current_users):
    receiver = st.selectbox("选择接收人", list(set(frineds["data"]["users"]) - set(current_users)))
    city = st.text_input("请准确填写城市，如'成都'")
    if st.button("确认修改", use_container_width=True, type="primary"):
        if reset_receiver(plugin_name="weather", args={"name": receiver, "city": city}).get("msg",
                                                                                            "") == "ok":
            st.rerun()


with weather_forecast:
    st.info("天气推送")
    weather_set = get_cron_setting("weather")
    if weather_set["users"]:
        if st.button("天气推送：" + weather_set["cron"]):
            update_cron_spec(weather_plugin_name, weather_set["cron"])
        st.title("已绑定账号")
        df = pd.DataFrame(weather_set["users"])
        df["check"] = True
        st.data_editor(df, key="friends_data_editor", column_order=[
            "name",
            "city",
            "check"
        ], column_config={
            "name": st.column_config.TextColumn("微信昵称(或备注)"),
            "city": st.column_config.TextColumn("关联查询城市"),
            "check": st.column_config.CheckboxColumn("已绑定"),
        }, on_change=weather_forecast_onlisten, hide_index=True, use_container_width=True)
    if st.button("新增接收人", type="primary", use_container_width=True):
        update_weather_receiver(df["name"].tolist() if len(df) != 0 else [])


@st.experimental_dialog("修改热点新闻配置")
def update_news_receiver():
    users = frineds["data"]["users"]
    gs = frineds["data"]['groups']
    receiver = st.multiselect("选择接收用户", users if users else [], default=news_set["users"])
    groups = st.multiselect("选择接收用户", gs if gs else [], default=news_set["groups"])
    if st.button("确认修改", use_container_width=True, type="primary"):
        if reset_receiver(plugin_name="news", args={"users": receiver, "groups": groups}).get("msg",
                                                                                              "") == "ok":
            st.rerun()


@st.experimental_dialog("RSS订阅源设置")
def rss(source, top_n):
    if source != "":
        # 关闭操作
        res = reset_rss("", top_n)
        st.rerun()
    new_source = st.text_input("RSS源地址", value=source)
    new_top_n = st.number_input("最多接受的消息量", min_value=5, max_value=50, step=1, value=top_n)
    if st.button("订阅", disabled=(new_source == source), type="primary",
                 use_container_width=True):
        res = reset_rss(new_source if new_source != source else source,
                        new_top_n if new_top_n != top_n else top_n)
        if res["msg"] == "ok":
            st.info("修改成功")
            st.rerun()
        else:
            st.warning(res["msg"])


def change_status_session():
    st.session_state["rss_toggle"] = True


with hot_search:
    st.header("实时订阅消息推送")
    news_set = get_cron_setting("news")
    if news_set["users"] or news_set["groups"]:
        if st.button("实时热点新闻推送：" + news_set["cron"]):
            update_cron_spec(news_plugin_name, news_set["cron"])

        st.toggle("使用RSS订阅", value=(news_set["rss_source"] != ""), on_change=change_status_session)
        if st.session_state.get("rss_toggle", False):
            del st.session_state["rss_toggle"]
            rss(news_set["rss_source"], news_set["top_n"])
        st.multiselect("当前已开启用户", options=news_set["users"], default=news_set["users"], disabled=True)
        st.multiselect("当前已开启群组", options=news_set["groups"], default=news_set["groups"], disabled=True)
    if st.button("修改配置", type="primary", use_container_width=True):
        update_news_receiver()
