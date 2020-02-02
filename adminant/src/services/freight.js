import request from '@/utils/request';
export async function listRule() {
  return request('/admin/frieght/list');
}
