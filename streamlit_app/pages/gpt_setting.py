import streamlit as st
from apis import *
res = check_login()
if res["code"] == 511:
    st.switch_page("home.py")
st.set_page_config(page_title='聊天模型配置', page_icon='🔩', layout='wide',
                   initial_sidebar_state="expanded")

# Title
st.title('✈️ 聊天模型管理')
# 提供修改选项
models = get_models()
text_model = models.get("data", {}).get("text_model", "---")
drawing_model = models.get("data", {}).get("drawing_model", "---")
# 展示当前使用的模型
text_c, picture_c = st.columns(2)
text_c.metric("文本模型", text_model, "using")
picture_c.metric("图片模型", drawing_model, "using")

choice_models = text_models()
# 提供修改聊天模型
with st.popover("修改聊天模型"):
    text_model = st.selectbox("修改聊天模型", options=choice_models, index=choice_models.index(text_model))
    if st.button("修改"):
        reset_models(text_model)
        st.rerun()
