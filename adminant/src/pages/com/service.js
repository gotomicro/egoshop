import request from '@/utils/request';

export async function fetchListRule(params) {
  return request('/admin/com/list', {
    params,
  });
}


export async function fetchOneRule(id) {
  return request('/admin/com/one/'+id, {});
}

// 获取文章内容
export async function contentRule(id) {
  return request('/admin/com/content/'+id, {});
}

export async function onSaleRule(params) {
  return request('/admin/com/onSale', {
    method: 'POST',
    data: { ...params},
  });
}

export async function offSaleRule(params) {
  return request('/admin/com/offSale', {
    method: 'POST',
    data: { ...params},
  });
}


export async function removeRule(params) {
  return request('/admin/com/remove', {
    method: 'POST',
    data: { ...params, method: 'delete' },
  });
}
export async function createRule(params) {
  return request('/admin/com/create', {
    method: 'POST',
    data: { ...params, method: 'post' },
  });
}
export async function updateRule(params) {
  return request('/admin/com/update', {
    method: 'POST',
    data: { ...params, method: 'update' },
  });
}

// com spec info


export async function comspecListRule(params) {
  return request('/admin/comspec/list', {
    params,
  });
}

export async function comspecCreateRule(params) {
  return request('/admin/comspec/create', {
    method: 'POST',
    data: { ...params },
  });
}


export async function comspecValueCreateRule(params) {
  return request('/admin/comspec/valueCreate', {
    method: 'POST',
    data: { ...params},
  });
}
