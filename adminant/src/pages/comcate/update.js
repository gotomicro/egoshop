import React, { Component } from "react";
import PageHeaderWrapper from "@/components/pageHeaderWrapper";
import { Card, message, Spin } from "antd";
import { connect } from "dva";
import { Form,  Input, Button, TreeSelect } from "antd";
import UploadImage from "@/components/uploadImage";
import router from "umi/router";
import Arr from "@/utils/array";

const FormItem = Form.Item;

@Form.create()
@connect(({ comcateModel, loading }) => ({
  cateTree: comcateModel.cateTree,
  listData: comcateModel.listData,
  infoData: comcateModel.infoData,
  loading: loading.effects["comcateModel/info"],

}))
export default class Index extends Component {
    componentDidMount() {
        const { location: { query: { id } }, dispatch } = this.props;
        dispatch({
            type: "comcateModel/info",
            payload: { id }
        });
        dispatch({
          type: 'comcateModel/fetchList',
        });
    }

    handleSubmit = (e) => {
        e.preventDefault();
        this.props.form.validateFields((err, values) => {
          console.log("fuck",err)
            if (!err) {
                const { location: { query: { id } }, dispatch } = this.props;
                dispatch({
                    type: "comcateModel/update",
                    payload: {
                        ...values,
                        ...{ id:Number(id) }
                    },
                    callback: (response) => {
                        if (response.code === 0) {
                            message.success("修改成功")
                            router.goBack();
                        }else{
                            message.error(response.msg)
                        }
                    }
                });
            }
        });
    };

    render() {
        const { form, listData,cateTree,infoData, loading } = this.props;
        const { name, pid, icon } = infoData || {};
        const { getFieldDecorator } = form;
        const formItemLayout = {
            labelCol: {
                xs: { span: 24 },
                sm: { span: 6 }
            },
            wrapperCol: {
                xs: { span: 24 },
                sm: { span: 18 }
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
                    offset: 6
                }
            }
        };
        let hasChild = false;
        listData.list.map((e) => {
            (hasChild === false && e.pid === infoData.id) ? hasChild = true : null;
        });

        // 最多3级 去掉自己
        let treeData = [];
        cateTree.forEach(function(item) {
            if (item.id !== infoData.id) {
                treeData.push({
                    title: item.name,
                    value: `${item.id}`,
                    key: `${item.id}`,
                    children: typeof item["children"] === "undefined" ? [] : ((item) => {
                        let _data = [];
                        item.children.map((sub) => {
                            if (sub.id !== infoData.id) {
                                _data.push({
                                    title: sub.name,
                                    value: `${sub.id}`,
                                    key: `${sub.id}`
                                });
                            }
                        });
                        return _data;
                    })(item)
                });
            }
        });
        return (
            <PageHeaderWrapper hiddenBreadcrumb={true}>
                <Card bordered={false}>
                    <Spin size="large" spinning={loading}>
                        <Form onSubmit={this.handleSubmit} style={{ maxWidth: 600 }}>
                            <FormItem
                                label="分类名称"
                                {...formItemLayout}
                            >
                                {getFieldDecorator("name", {
                                    initialValue: name,
                                    rules: [{
                                        required: true,
                                        message: "请输入分类名称"
                                    }]
                                })(
                                    <Input
                                        placeholder='请输入分类名称'
                                        style={{ width: "100%" }}
                                    />
                                )}
                            </FormItem>
                            {hasChild === false ? <FormItem
                                label="上级分类"
                                help="如不选择，则默认为一级分类"
                                {...formItemLayout}
                            >
                                {getFieldDecorator("pid", {
                                    initialValue: pid === 0 ? 0 : pid
                                })(
                                    <TreeSelect
                                        treeData={treeData}
                                        showSearch
                                        dropdownStyle={{ maxHeight: 400, overflow: "auto" }}
                                        placeholder="请输入分类名称"
                                        allowClear
                                        treeDefaultExpandAll
                                    />
                                )}
                            </FormItem> : null}
                            <FormItem
                                {...formItemLayout}
                                extra="分类展示图，建议尺寸：140*140 像素"
                                label="上传分类图"
                            >
                                {getFieldDecorator("icon", {
                                    initialValue: icon,
                                    rules: [{
                                        message: "请上传分类图"
                                    }],
                                    valuePropName: "url"
                                })(
                                    <UploadImage />
                                )}
                            </FormItem>
                            <FormItem {...tailFormItemLayout}>
                                <Button
                                    type="primary"
                                    htmlType="submit"
                                    style={{
                                        marginRight: 10
                                    }}
                                >
                                    保存
                                </Button>
                                <Button
                                    onClick={() => {
                                        router.goBack();
                                    }}
                                >
                                    返回
                                </Button>
                            </FormItem>
                        </Form>
                    </Spin>
                </Card>
            </PageHeaderWrapper>
        );
    }

}
