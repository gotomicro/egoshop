import {
  Button,
  Card,
  Col,
  Form,
  Icon,
  Input,
  Row,
  Select,
  message,
  Divider,
  Popconfirm,
  Switch,
} from 'antd';
import React, { Component, Fragment } from 'react';
import { PageHeaderWrapper } from '@ant-design/pro-layout';
import { connect } from 'dva';
import moment from 'moment';
import StandardTable from '@/components/SaberStandardTable';
import styles from './style.less';
import router from "umi/router";

const FormItem = Form.Item;
const { Option } = Select;
import Image from "@/components/SaberImage";

const getValue = obj =>
  Object.keys(obj)
    .map(key => obj[key])
    .join(',');


/* eslint react/no-multi-comp:0 */
@connect(({ comModel, loading }) => ({
  listData: comModel.listData,
  loading: loading.models.feedModel,
}))
class TableList extends Component {
  state = {
    modalVisible: false,
    updateModalVisible: false,
    expandForm: false,
    selectedRows: [],
    formValues: {},
    updateInitialValues: {},
  };

  columns = [
    {
      title: '编号',
      dataIndex: 'id',
    },
    {
      title: "商品图",
      dataIndex: "cover",
      width: 50,
      render: (e) => (
        <Image
          type='goods'
          src={e}
          style={{ width: 50, height: 50 }}
        />
      )
    }, {
      title: "商品标题",
      dataIndex: "title",
      width: 230
    }, {
      title: "价格（元）",
      dataIndex: "price",
      width: 120
    }, {
      title: "销量",
      dataIndex: "baseSaleNum",
      width: 80
    }, {
      title: "库存",
      dataIndex: "stock",
      width: 80
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      sorter: true,
      render: val => <span>{moment(val).format('YYYY-MM-DD HH:mm:ss')}</span>,
    },
    {
      title: "上架状态",
      dataIndex: "isOnSale",
      render: (text, record) => <Switch
        defaultChecked={text}
        onChange={async (checked) => {
          const { dispatch } = this.props;
          if (checked) {
            dispatch({
              type: 'comModel/onSale',
              payload: { ids: [record.id] },
            });
          } else {
            dispatch({
              type: 'comModel/offSale',
              payload: { ids: [record.id] },
            });
          }
        }}
      />
    },
    {
      title: '操作',
      render: (text, record) => (
        <Fragment>
          <a
            onClick={() => {
              router.push({
                pathname: `/com/update`,
                search: `?id=${record.id}`,
                state: {
                  record
                }
              });
            }}
          >
            修改
          </a>
          <Divider type="vertical" />
          <a
            onClick={() => {
              router.push({
                pathname: `/com/editor`,
                search: `?id=${record.id}`,
                state: {
                  record
                }
              });
            }}
          >
            内容
          </a>
          <Divider type="vertical" />
          <Popconfirm title="是否要删除此行？" onConfirm={() => this.handleRemove(record.id)}>
            <a>删除</a>
          </Popconfirm>
        </Fragment>
      ),
    },
  ];

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({
      type: 'comModel/fetchList',
    });
  }

  handleStandardTableChange = (pagination, filtersArg, sorter) => {
    const { dispatch } = this.props;
    const { formValues } = this.state;
    const filters = Object.keys(filtersArg).reduce((obj, key) => {
      const newObj = { ...obj };
      newObj[key] = getValue(filtersArg[key]);
      return newObj;
    }, {});
    const params = {
      currentPage: pagination.current,
      pageSize: pagination.pageSize,
      ...formValues,
      ...filters,
    };

    if (sorter.field) {
      params.sorter = `${sorter.field}_${sorter.order}`;
    }

    dispatch({
      type: 'comModel/fetchList',
      payload: params,
    });
  };

  handleFormReset = () => {
    const { form, dispatch } = this.props;
    form.resetFields();
    this.setState({
      formValues: {},
    });
    dispatch({
      type: 'comModel/fetchList',
      payload: {},
    });
  };

  handleSearch = e => {
    e.preventDefault();
    const { dispatch, form } = this.props;
    form.validateFields((err, fieldsValue) => {
      if (err) return;
      const values = {
        ...fieldsValue,
        updatedAt: fieldsValue.updatedAt && fieldsValue.updatedAt.valueOf(),
      };
      this.setState({
        formValues: values,
      });
      dispatch({
        type: 'comModel/fetchList',
        payload: values,
      });
    });
  };

  handleRemove = id => {
    const { dispatch } = this.props;
    dispatch({
      type: 'comModel/remove',
      payload: {
        id,
      },
    });
    message.success('添加成功');
  };

  renderSimpleForm() {
    const { form } = this.props;
    const { getFieldDecorator } = form;
    return (
      <Form onSubmit={this.handleSearch} layout="inline">
        <Row
          gutter={{
            md: 8,
            lg: 24,
            xl: 48,
          }}
        >
          <Col md={6} sm={24}>
            <FormItem label="名称">
              {getFieldDecorator('name')(<Input placeholder="请输入名称" />)}
            </FormItem>
          </Col>
          <Col md={6} sm={24}>
            <FormItem label="繁育人">
              {getFieldDecorator('author')(<Input placeholder="请输入繁育人" />)}
            </FormItem>
          </Col>
          <Col md={6} sm={24}>
            <FormItem label="地址">
              {getFieldDecorator('address')(<Input placeholder="请输入地址" />)}
            </FormItem>
          </Col>
          <Col md={6} sm={24}>
            <span className={styles.submitButtons}>
              <Button type="primary" htmlType="submit">
                查询
              </Button>
              <Button
                style={{
                  marginLeft: 8,
                }}
                onClick={this.handleFormReset}
              >
                重置
              </Button>
            </span>
          </Col>
        </Row>
      </Form>
    );
  }

  render() {
    const {
      listData,
      loading,
    } = this.props;
    return (
      <PageHeaderWrapper>
        <Card bordered={false}>
          <div className={styles.tableList}>
            <div className={styles.tableListForm}>{this.renderSimpleForm()}</div>
            <div className={styles.tableListOperator}>
              <Button icon="plus" type="primary" onClick={() => {
                router.push({
                  pathname: `/com/create`,
                });
              }}>
                新建
              </Button>
            </div>
            <StandardTable
              loading={loading}
              data={listData}
              columns={this.columns}
              onSelectRow={this.handleSelectRows}
              onChange={this.handleStandardTableChange}
            />
          </div>
        </Card>
      </PageHeaderWrapper>
    );
  }
}

export default Form.create()(TableList);
