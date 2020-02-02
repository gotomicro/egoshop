import { fetchListRule } from './service';

const Model = {
  namespace: 'usersModel',
  state: {
    listData: {
      list: [],
      pagination: {},
    },
  },
  effects: {
    *fetchList({ payload }, { call, put }) {
      const response = yield call(fetchListRule, payload);
      yield put({
        type: '_fetchList',
        payload: response.data,
      });
    },
  },
  reducers: {
    _fetchList(state, action) {
      return {
        ...state,
        listData: action.payload
      };
    },
  },
};
export default Model;
