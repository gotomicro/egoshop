import { fetchListRule, fetchOneRule,contentRule, onSaleRule,offSaleRule,createRule, removeRule, updateRule,
  // comspec info
  comspecListRule,
  comspecCreateRule,
  comspecValueCreateRule,
} from './service';

const Model = {
  namespace: 'comModel',
  state: {
    listData: {
      list: [],
      pagination: {},
    },
    oneData: {
      gallery: [],
      skuList: [],
      saleTime: "",
    },
    contentData: {
    },
    comspecListData: {
      list: [],
      pagenation: {},
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

    *fetchOne({ payload, callback }, { call, put }) {
      const response = yield call(fetchOneRule, payload.id);
      yield put({
        type: '_fetchOne',
        payload: response.data,
      });
      if (callback) callback(response);
    },

    *content({ payload }, { call, put }) {
      const response = yield call(contentRule, payload.id);
      yield put({
        type: '_content',
        payload: response.data,
      });
    },

    *onSale({ payload }, { call, put }) {
      const response = yield call(onSaleRule, payload);
      yield put({
        type: '_fetchList',
        payload: response.data,
      });
    },
    *offSale({ payload }, { call, put }) {
      const response = yield call(offSaleRule, payload);
      yield put({
        type: '_fetchList',
        payload: response.data,
      });
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
      if (callback) callback(response);
    },

    *update({ payload, callback }, { call, put }) {
      const response = yield call(updateRule, payload);
      yield put({
        type: '_fetchList',
        payload: response.data,
      });
      if (callback) callback(response);
    },
    *comspecList({ payload, callback }, { call, put }) {
      const response = yield call(comspecListRule, payload);
      yield put({
        type: '_comspecList',
        payload: response.data,
      });
      if (callback) callback(response);
    },
    *comspecCreate({ payload, callback }, { call, put }) {
      const response = yield call(comspecCreateRule, payload);
      if (callback) callback(response);
    },
    *comspecValueCreate({ payload, callback }, { call, put }) {
      const response = yield call(comspecValueCreateRule, payload);
      if (callback) callback(response);
    },

  },
  reducers: {
    _fetchList(state, action) {
      return {
        ...state,
        listData: action.payload
      };
    },
    _fetchOne(state, action) {
      return {
        ...state,
        oneData: action.payload
      };
    },
    // 文档内容
    _content(state, action) {
      return {
        ...state,
        contentData: action.payload
      };
    },
    _comspecList(state, action) {
      return {
        ...state,
        comspecListData: action.payload
      };
    },
  },
};
export default Model;
