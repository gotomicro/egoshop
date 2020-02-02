import request from '@/utils/request';

export async function fetchListRule(params) {
  return request('/admin/comcate/list', {
    params,
  });
}

export async function infoRule(params) {
  return request('/admin/comcate/info', {
    params,
  });
}
export async function removeRule(params) {
  return request('/admin/comcate/remove', {
    method: 'POST',
    data: { ...params, method: 'delete' },
  });
}
export async function createRule(params) {
  return request('/admin/comcate/create', {
    method: 'POST',
    data: { ...params, method: 'post' },
  });
}
export async function updateRule(params) {
  return request('/admin/comcate/update', {
    method: 'POST',
    data: { ...params, method: 'update' },
  });
}

