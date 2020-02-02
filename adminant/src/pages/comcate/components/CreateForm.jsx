import { Form, Input, Modal } from 'antd';
import React from 'react';
import UploadImage from "@/components/uploadImage";

const FormItem = Form.Item;

const CreateForm = props => {
  const { modalVisible, form, handleAdd, handleModalVisible } = props;

  const okHandle = () => {
    form.validateFields((err, fieldsValue) => {
      if (err) return;
      form.resetFields();
      handleAdd(fieldsValue);
    });
  };

  return (
    <Modal
      destroyOnClose
      title="新建分类"
      visible={modalVisible}
      onOk={okHandle}
      onCancel={() => handleModalVisible()}
    >
      <FormItem
        labelCol={{
          span: 5,
        }}
        wrapperCol={{
          span: 15,
        }}
        label="名称"
      >
        {form.getFieldDecorator('name', {
          rules: [
            {
              required: true,
              message: '请输入至少两个字符的规则描述！',
              min: 2,
            },
          ],
        })(<Input placeholder="请输入" />)}
      </FormItem>
      <FormItem
        labelCol={{
          span: 5,
        }}
        wrapperCol={{
          span: 15,
        }}
        extra="分类展示图，建议尺寸：140*140 像素"
        label="上传分类图"
      >
        {form.getFieldDecorator("icon", {
          rules: [{
            message: "请上传分类图"
          }],
          valuePropName: "url"
        })(
          <UploadImage />
        )}
      </FormItem>
    </Modal>
  );
};

export default Form.create()(CreateForm);
