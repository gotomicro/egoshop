import fa from "@/utils/fa";
import request from '@/utils/request';

export default {
    async list(data = {}) {
        return fa.request({
            url: `/admin/image/list`,
            method: "GET",
            data
        });
    },
  async add(params) {
    return request('/admin/image/add',{
      method: 'POST',
      data: { ...params},
    });
  }

};

