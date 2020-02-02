import { routerRedux } from 'dva/router';
import { stringify } from 'querystring';
import { fakeAccountLogin, getFakeCaptcha,accountLogout } from '@/services/login';
import { setAuthority } from '@/utils/authority';
import { getPageQuery } from '@/utils/utils';
import { message } from "antd";


const Model = {
  namespace: 'login',
  state: {
    code: 0,
  },
  effects: {
    *login({ payload }, { call, put }) {
      const response = yield call(fakeAccountLogin, payload);
      yield put({
        type: 'changeLoginStatus',
        payload: response,
      }); // Login successfully
      if (response.code === 0) {
        const params = getPageQuery();
        let { redirect } = params;
        if (redirect) {
          window.location.href = redirect;
        }else {
          window.location.href = "/"
        }
      }else {
        message.error("登录失败，请检查用户名或密码是否正确！")
      }
    },

    *getCaptcha({ payload }, { call }) {
      yield call(getFakeCaptcha, payload);
    },

    *logout(_, { call, put  }) {
      const { redirect } = getPageQuery(); // redirect

      const response = yield call(accountLogout);
      // Login successfully
      console.log('logout response', response);

      if (response.code === 0) {
        yield put({
          type: 'changeLoginStatus',
          payload: {
            status: false,
            data: {
              currentAuthority: 'guest',
            },
          },
        });
        yield put(
          routerRedux.push({
            pathname: '/user/login',
            search: stringify({
              redirect: window.location.href,
            }),
          })
        );
      }

    },
  },
  reducers: {
    changeLoginStatus(state, { payload }) {
      if (payload.code === 0) {
        setAuthority(payload.data.currentAuthority);
      }else {
        setAuthority({
          currentAuthority: 'guest',
        });
      }
      return { ...state, code: payload.code, type: payload.type };
    },
  },
};
export default Model;
