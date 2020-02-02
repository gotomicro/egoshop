import { Button, DatePicker, Form, Input, Modal, Radio, Select, Steps } from 'antd';
import React, { Component } from 'react';
import UploadImage from "@/components/uploadImage";

const FormItem = Form.Item;
const { Option } = Select;
const RadioGroup = Radio.Group;

class UpdateForm extends Component {
  static defaultProps = {
    handleUpdate: () => {},
    handleUpdateModalVisible: () => {},
    values: {},
  };

  formLayout = {
    labelCol: {
      span: 7,
    },
    wrapperCol: {
      span: 13,
    },
  };

  constructor(props) {
    super(props);
    this.state = {
      formVals: {
        ...props.values
      },
    };
  }

  okHandle = () => {
    const { form,handleUpdate } = this.props;
    const {formVals} = this.state
    form.validateFields((err, fieldsValue) => {
      if (err) return;
      form.resetFields();
      console.log("formvals",formVals)
      const mergeValue = { ...formVals, ...fieldsValue };
      handleUpdate(mergeValue);
    });
  };

  renderContent = (formVals) => {
    console.log("formVals content",formVals)
    const { form } = this.props;
      return [
        <FormItem key="name" {...this.formLayout} label="名称">
          {form.getFieldDecorator('name', {
            initialValue: formVals.name,
          })(
            <Input
              placeholder="请输入名称"
            />
          )}
        </FormItem>,
        <FormItem
          key="icon"
          {...this.formItemLayout}
          extra="分类展示图，建议尺寸：140*140 像素"
          label="上传分类图"
        >
          {form.getFieldDecorator("icon", {
            initialValue: formVals.icon,
            rules: [{
              message: "请上传分类图"
            }],
            valuePropName: "url"
          })(
            <UploadImage />
          )}
        </FormItem>
      ];
  };

  render() {
    const { updateModalVisible, handleUpdateModalVisible, values } = this.props;
    const { formVals } = this.state;
    return (
      <Modal
        width={640}
        bodyStyle={{
          padding: '32px 40px 48px',
        }}
        destroyOnClose
        title="修改分类信息"
        visible={updateModalVisible}
        onOk={this.okHandle}
        onCancel={() => handleUpdateModalVisible(false, values)}
        afterClose={() => handleUpdateModalVisible()}
      >
        {this.renderContent( formVals)}
      </Modal>
    );
  }
}

export default Form.create()(UpdateForm);
