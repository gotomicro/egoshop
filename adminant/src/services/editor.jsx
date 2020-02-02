import request from '@/utils/request';


export async function contentSaveRule(params) {
  return request('/admin/editor/contentSave', {
    method: 'POST',
    data: { ...params},
  });
}

export async function releaseRule(params) {
  return request('/admin/editor/release', {
    method: 'POST',
    data: { ...params},
  });
}


// 上传图片
export async function uploadRule(params) {
  return request('/admin/editor/upload', {
    method: 'POST',
    data: { ...params},
  });
}

