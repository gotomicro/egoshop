import { listRule } from "@/services/freight";

export default {
  namespace: "freightModel",
  state: {
    listData: {
      list: [],
      pagination: {}
    }
  },
  effects: {
    * list({ payload, callback }, { call, put }) {
      const response = yield call(listRule, payload);
      yield put({
        type: "_list",
        payload: response.data
      });
      if (callback) callback(response);
    }
  },
  reducers: {
    _list(state, action) {
      return {
        ...state,
        list: action.payload
      };
    }
  }
};
