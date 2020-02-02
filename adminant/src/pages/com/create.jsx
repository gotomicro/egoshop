import React, { Component } from "react";
import { connect } from "dva";
import {
  Form,
  Button,
  Modal,
  message,
  Card,
  Spin,
  Input,
  Select,
  TreeSelect,
  DatePicker
} from "antd";
import PageHeaderWrapper from "@/components/pageHeaderWrapper";
import PhotoGallery from "@/components/public/photoGallery";
import styles from "./edit.css";
import { UploadGroupImage } from "@/components/uploadImage";
import router from "umi/router";
import Sku from "./createSku";
import Freight from "./createFreight";

import moment from "moment"
import Antd from "@/utils/antd";

const { Option } = Select;
const FormItem = Form.Item;

@Form.create()
@connect(({ comcateModel, loading }) => ({
  cateList: comcateModel.listData.list,
  cateTree: comcateModel.cateTree,
  cateListLoading: loading.effects["comcateModel/fetchList"],
}))

class GoodsEdit extends Component {
  static defaultProps = {
    cateListLoading: true,
  };
  state = {
    photoGalleryVisible: false,
    photoGalleryOnOk: (e) => {
    },
    previewVisible: false,
    previewImage: "",
    save: true,
  }
  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({
      type: "comcateModel/fetchList"
    });
    let value = JSON.parse(localStorage.getItem("xixun-com-create")) || {}
    this.props.form.setFieldsValue(
      value.saleTime ? {
        ...value,
        saleTime: moment(value.saleTime)
      } : value
    )
  }

  componentWillUnmount(){
    if(this.state.save){
      let value = this.props.form.getFieldsValue()
      localStorage.setItem("xixun-com-create", JSON.stringify(value))
    }else {
      localStorage.removeItem("xixun-com-create")
    }
  }

  handleSubmit = (e) => {
    e.preventDefault();
    this.props.form.validateFieldsAndScroll(async (err, values) => {
      if (!err) {
        const { dispatch } = this.props;
        dispatch({
          type: "comModel/create",
          payload: {
            ...values,
            cid: parseInt(values.cid),
          },
          callback: (e) => {
            if (e.code === 0) {
              this.setState({
                save: false
              },()=>{
                message.success("添加成功");
                router.goBack();
              })
            } else {
              message.warn(e.msg);
            }
          }
        });
      }
    });
  };


  render() {
    const { photoGalleryVisible, previewVisible, previewImage } = this.state;
    const { form, cateListLoading,cateList,cateTree } = this.props;
    const cateTreeSelect = Antd.treeData(cateTree);
    const { getFieldDecorator, getFieldValue, setFieldsValue } = form;

    const selectOptions = cateList.map(d => <Option key={d.id}>{d.name}</Option>);

    return (
      <PageHeaderWrapper hiddenBreadcrumb={true}>
        <Card bordered={false}>
          <Spin size="large" spinning={cateListLoading}>
            <Form onSubmit={this.handleSubmit} style={{ width: 1000 }}>
              <div className={styles.item}>
                <h3>基本信息</h3>
                <FormItem
                  {...formItemLayout}
                  label='商品图'
                >
                  {getFieldDecorator("gallery", {
                    rules: [{ required: true, message: "请选择商品图" }],
                    valuePropName: "url",
                  })(
                    <UploadGroupImage
                      onClick={(onChange, values) => {
                        values = values ? values : [];
                        this.openPhotoGallery({
                          photoGalleryOnOk: (e) => {
                            onChange([...values, ...e]);
                          }
                        });
                      }}
                      preview={(previewImage) => {
                        this.openPreviewModal({
                          previewImage
                        });
                      }}
                    />
                  )}
                </FormItem>
                <FormItem
                  {...formItemLayout}
                  label='商品名称'
                >
                  {getFieldDecorator("title", {
                    rules: [{ required: true, message: "请输入商品名称" }],
                  })(
                    <Input
                      placeholder="请输入商品名称"
                    />
                  )}
                </FormItem>
                <FormItem
                  {...formItemLayout}
                  label='商品副名称'
                >
                  {getFieldDecorator("subTitle", {
                    rules: [{ required: true, message: "商品副名称" }],
                  })(
                    <Input
                      placeholder="商品副名称"
                    />
                  )}
                </FormItem>
                <FormItem
                  {...formItemLayout}
                  label='商品分类'
                >
                  {getFieldDecorator("cids", {
                    rules: [{ required: true, message: "请选择商品分类" }]
                  })(
                    <TreeSelect
                      treeData={cateTreeSelect}
                      showSearonChangech
                      dropdownStyle={{ maxHeight: 400, overflow: "auto" }}
                      placeholder="请选择商品分类"
                      allowClear
                      multiple
                      treeDefaultExpandAll
                      onChange={(value) => {
                        setFieldsValue({
                          cids: value
                        });
                      }}
                    />
                  )}
                  <a
                    onClick={() => {
                      router.push("/comcate/create");
                    }}
                  >
                    新增分类
                  </a>
                </FormItem>
              </div>
              <div className={styles.item}>
                <h3>型号价格</h3>
                <FormItem {...formItemLayout}>
                  {getFieldDecorator("skuList", {
                    rules: [{
                      validator: Sku.validator,
                      required: true
                    }],
                  })(<Sku form={form} />)}
                </FormItem>
              </div>
              <div className={styles.item}>
                <h3>运费其他</h3>
                <FormItem {...formItemLayout} label={"运费"}>
                  {getFieldDecorator("freight", {
                    rules: [{
                      required: true
                    }],
                  })(<Freight />)}
                </FormItem>
                <FormItem {...formItemLayout} label={"开售时间"}>
                  {getFieldDecorator("saleTime", {
                    rules: [{
                      required: true
                    }],
                  })(
                    <DatePicker
                      showTime
                      format="YYYY-MM-DD HH:mm:ss"
                      placeholder="选择时间"
                      style={{ marginRight: 15 }}
                    />
                  )}
                </FormItem>
              </div>
              <FormItem {...tailFormItemLayout}>
                <Button
                  type="primary"
                  htmlType="submit"
                  style={{
                    marginRight: 10
                  }}
                >
                  保存添加
                </Button>
                <Button
                  onClick={()=>{
                    router.goBack()
                  }}
                >
                  返回
                </Button>
              </FormItem>
            </Form>
            <PhotoGallery
              visible={photoGalleryVisible}
              onCancel={this.onCancelPhotoGallery}
              onOk={this.onOkPhotoGallery}
            />
            <Modal visible={previewVisible} footer={null} onCancel={this.previewCancel}>
              <img alt="example" style={{ width: "100%" }} src={previewImage} />
            </Modal>
          </Spin>
        </Card>
      </PageHeaderWrapper>
    );
  }


  openPhotoGallery = ({ photoGalleryOnOk }) => {
    this.setState({
      photoGalleryVisible: true,
      photoGalleryOnOk
    });
  };
  onCancelPhotoGallery = () => {
    this.setState({
      photoGalleryVisible: false
    });
  };
  onOkPhotoGallery = (e) => {
    this.state.photoGalleryOnOk(e);
    this.onCancelPhotoGallery();
  };
  previewCancel = () => {
    this.setState({
      previewVisible: false
    });
  };
  // : { previewImage: string }
  openPreviewModal = ({ previewImage }) => {
    this.setState({
      previewVisible: true,
      previewImage
    });
  };

}

const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 4 }
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 20 }
  }
};
const tailFormItemLayout = {
  wrapperCol: {
    xs: {
      span: 24,
      offset: 0
    },
    sm: {
      span: 16,
      offset: 4
    }
  }
};

export default GoodsEdit;
