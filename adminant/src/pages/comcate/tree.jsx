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
  Table,
} from 'antd';
import React, { Component, Fragment } from 'react';
import { PageHeaderWrapper } from '@ant-design/pro-layout';
import { connect } from 'dva';
import moment from 'moment';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import StandardTable from '@/components/SaberStandardTable';
import styles from './style.less';
import Image from "@/components/SaberImage";
import Arr from "@/utils/array";
import CategorySort from "@/components/goods/category/sort";
import { View } from "@/components/public/dom";
import router from "umi/router";

const FormItem = Form.Item;
const { Option } = Select;

const getValue = obj =>
  Object.keys(obj)
    .map(key => obj[key])
    .join(',');


/* eslint react/no-multi-comp:0 */
@connect(({ comcateModel, loading }) => ({
  listData: comcateModel.listData,
  cateTree: comcateModel.cateTree,
  loading: loading.models.comcateModel,
}))
class TableList extends Component {
  state = {
    modalVisible: false,
    updateModalVisible: false,
    expandForm: false,
    selectedRows: [],
    formValues: {},
    updateInitialValues: {},
    expandedRowKeys: [],
  };

  columns = [
    {
      title: '编号',
      dataIndex: 'id',
    },
    {
      title: '名称',
      dataIndex: 'name',
    },
    {
      title: '图片',
      dataIndex: 'icon',
      render: val => <Image src={val} style={{ height: "20px" }} />
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      sorter: true,
      render: val => <span>{moment(val).format('YYYY-MM-DD HH:mm:ss')}</span>,
    },
    {
      title: '更新时间',
      dataIndex: 'updatedAt',
      sorter: true,
      render: val => <span>{moment(val).format('YYYY-MM-DD HH:mm:ss')}</span>,
    },
    {
      title: '操作',
      render: (text, record) => (
        <Fragment>
          <a onClick={() => this.handleUpdateModalVisible(true, record)}>修改</a>
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
      type: 'comcateModel/fetchList'
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
      type: 'comcateModel/fetchList',
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
      type: 'comcateModel/fetchList',
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
        type: 'comcateModel/fetchList',
        payload: values,
      });
    });
  };

  handleModalVisible = flag => {
    this.setState({
      modalVisible: !!flag,
    });
  };

  handleUpdateModalVisible = (flag, record) => {
    this.setState({
      updateModalVisible: !!flag,
      updateInitialValues: record || {},
    });
  };

  handleAdd = fields => {
    console.log("fields",fields)
    const { dispatch } = this.props;
    dispatch({
      type: 'comcateModel/add',
      payload: {
        ...fields,
      },
    });
    message.success('添加成功');
    this.handleModalVisible();
  };

  handleUpdate = fields => {
    const { dispatch } = this.props;
    dispatch({
      type: 'comcateModel/update',
      payload: {
        ...fields,
      },
    });
    message.success('配置成功');
    this.handleUpdateModalVisible();
  };


  handleRemove = id => {
    const { dispatch } = this.props;
    dispatch({
      type: 'comcateModel/remove',
      payload: {
        id,
      },
    });
    message.success('添加成功');
    this.handleModalVisible();
  };

  render() {
    const { expandedRowKeys } = this.state;
    const { loading, listData, dispatch,cateTree } = this.props;

    console.log("cateTree",cateTree)

    const columns = [
      {
        title: "ID - 分类名称",
        dataIndex: "name",
        key: "name",
        render: (value, row) => {
          return <span><span style={{ marginRight: 15 }}><span
            style={{ color: "#ccc", fontSize: 10 }}>{row.id} - </span>{value}</span><span><Image
            src={row.icon} style={{ height: "20px" }} /></span></span>;
        }
      },

      {
        title: "操作",
        key: "operation",
        className: styles.column,
        render: (record) => <Fragment>
          <a
            onClick={() => {
              router.push({
                pathname: `/comcate/update`,
                search: `?id=${record.id}`,
                state: {
                  categoryData: record
                }
              });
            }}
          >
            编辑
          </a>
          <Divider type="vertical" />
          <Popconfirm
            title="确认删除？"
            okText="确认"
            cancelText="取消"
            onConfirm={() => {
              dispatch({
                type: 'comcateModel/remove',
                payload: {
                  id: record.id
                }
              });
            }}
          >
            <a>删除</a>
          </Popconfirm>
        </Fragment>
      }
    ];

    return (
      <PageHeaderWrapper>
        <Card bordered={false}>
          <View>
            <div className={styles.batchView}>
              <Button
                type='primary'
                onClick={() => {
                  router.push("/comcate/create");
                }}
              >
                添加分类
              </Button>
              <Button
                onClick={() => {
                  this.setState({
                    expandedRowKeys: Array.isArray(listData.list) ? listData.list.map((item) => item.id) : []
                  });
                }}
              >
                全部展开
              </Button>
              <Button
                onClick={() => {
                  this.setState({
                    expandedRowKeys: []
                  });
                }}
              >
                全部折叠
              </Button>
              <CategorySort dataSource={cateTree} dispatch={dispatch} />
            </div>
            <Table
              loading={loading}
              columns={columns}
              dataSource={cateTree}
              defaultExpandAllRows={true}
              rowKey={record => record.id}
              pagination={false}
              expandedRowKeys={expandedRowKeys}
              onExpand={(bool, row) => {
                if (bool) {
                  this.setState({
                    expandedRowKeys: [...expandedRowKeys, row.id]
                  });
                } else {
                  const index = expandedRowKeys.findIndex((e) => e === row.id);
                  const newArray = [...expandedRowKeys];
                  newArray.splice(index, 1);
                  this.setState({
                    expandedRowKeys: newArray
                  });
                }
              }}
            />
          </View>
        </Card>
      </PageHeaderWrapper>
    );
  }
}

export default Form.create()(TableList);
