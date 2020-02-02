import request from '@/utils/request';

export async function fetchListRule(params) {
  return request('/admin/users/list', {
    params,
  });
}
