import { contentSaveRule,uploadRule,releaseRule } from "@/services/editor";

const Feed = {
  namespace: 'editorModel',
  state: {
  },
  effects: {
    *release({ payload, callback }, { call, put }) {
      const response = yield call(releaseRule, payload);
      if (callback) callback(response);
    },
    *contentSave({ payload, callback }, { call, put }) {
      const response = yield call(contentSaveRule, payload);
      if (callback) callback(response);
    },
    // 上传图片
    *upload({ payload , callback }, { call, put }) {
      const response = yield call(uploadRule, payload);
      if (callback) callback(response);
    },
  },
  reducers: {
  },
};
export default Feed;
