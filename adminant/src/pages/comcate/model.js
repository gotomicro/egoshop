import { createRule,infoRule, fetchListRule, removeRule, updateRule } from './service';
import Arr from "@/utils/array";

const Model = {
  namespace: 'comcateModel',
  state: {
    listData: {
      list: [],
      pagination: {},
    },
    infoData: {},
    cateTree: [],
  },
  effects: {
    *fetchList({ payload,callback }, { call, put }) {
      const response = yield call(fetchListRule, payload);
      yield put({
        type: '_fetchList',
        payload: response.data,
      });
      if (callback) callback(response.data)
    },
    *info({ payload,callback }, { call, put }) {
      const response = yield call(infoRule, payload);
      yield put({
        type: '_info',
        payload: response.data,
      });
      if (callback) callback(response.data)
    },

    *create({ payload, callback }, { call, put }) {
      const response = yield call(createRule, payload);
      yield put({
        type: '_fetchList',
        payload: response.data,
      });
      if (callback) callback(response);
    },

    *remove({ payload, callback }, { call, put }) {
      const response = yield call(removeRule, payload);
      yield put({
        type: '_fetchList',
        payload: response.data,
      });
      if (callback) callback();
    },

    *update({ payload, callback }, { call, put }) {
      const response = yield call(updateRule, payload);
      yield put({
        type: '_fetchList',
        payload: response.data,
      });
      if (callback) callback(response);
    },
  },
  reducers: {
    _fetchList(state, action) {
      const cateTree = Arr.toTree(action.payload.list);
      return {
        ...state,
        listData: action.payload,
        cateTree: cateTree
      };
    },
    _info(state, action) {
      return {
        ...state,
        infoData: action.payload,
      };
    },
  },
};
export default Model;
