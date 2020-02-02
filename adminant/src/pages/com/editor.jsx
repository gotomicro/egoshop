import React, { Component } from 'react';
import style from './editor.less';
import {Card, Col, Icon, Input, message, Modal, Row, Select, Tabs, Tree, Button} from 'antd';
const Search = Input.Search;
import { ShowPanel,EditPanel } from "@/components/SaberEditor";
import UtilTime from "@/utils/time"
import { connect } from 'dva';

function nameFomater(name) {
  if (name === "") {
    return '暂未记录';
  }
  return name;
}

@connect(({ comModel, loading }) => ({
  contentData: comModel.contentData,
  loading: loading.models.comModel,
}))
export default class Editor extends Component {
    state = {
      isEdit: false,
      id: 0,
      type: 2,
    };

    constructor(props) {
      super(props);
    }

  componentDidMount() {
    const { dispatch, location: { query: { id } } } = this.props;
    if (id === undefined || id === "" || id === "0") {
      message.error("doc is is empty")
      return
    }
    // 设置文档id
    this.setState({id:parseInt(id)})
    dispatch({
      type: 'comModel/content',
      payload: {
        id,
      },
    });
  }

  callbackContent = () => {
    const { dispatch } = this.props;
    const {id} = this.state;
    dispatch({
      type: 'comModel/content',
      payload: {
        id:id,
      },
    });
  }

  releaseContent = () => {
    const { dispatch } = this.props;
    const {id,type} = this.state;
    dispatch({
      type: 'editorModel/release',
      payload: {
        id:id,
        type: type,
      },
    });
  }

  render() {
    const { isEdit,type } = this.state
    const { contentData } = this.props;
    return (
      <div>
          <div className={style.docShow}>
            <div className={style.docContainer}>
              <Row className={style.docHeader}>
                <div className={style.docHeaderLeft}>
                  <div>标题：{contentData.name}, 最近编辑时间：{`${UtilTime.relativeTime(contentData.updatedAt)}，`} 最近编辑者：{nameFomater(contentData.updatedBy)}</div>
                </div>
                <div className={style.docHeaderRight}>
                  {isEdit === true && (
                    <div className={style.docBtnEditor} onClick={() => {
                      this.setState({
                        isEdit: !isEdit,
                      });
                    }}><i className="fa fa-edit" name="edit" />阅读</div>
                  )}
                  {isEdit === false && (
                    <div className={style.docBtnEditor} onClick={() => {
                      this.setState({
                        isEdit: !isEdit,
                      });
                    }}><i className="fa fa-edit" name="edit" />编辑</div>
                  )}
                  <div className={style.docBtnEditor} onClick={() => {
                    this.releaseContent()
                  }}><i className="fa fa-edit" name="edit" />发布</div>
                  <div className={style.toolbar}>
                    <div className={style.toolsItem} title="分享" onClick={() => {
                    }}>
                      <i className="fa fa-arrows-alt" name="edit" />
                    </div>
                    <div className={style.toolsItem} title="查看标签">
                      <i className="fa fa-edit" name="edit" />
                    </div>
                    <div className={style.toolsItem} title="更多" widget-toggle="false">
                      <i className="fa fa-edit" name="edit" />
                    </div>
                  </div>
                </div>
              </Row>
              {isEdit === true && (
                <Row className={style.docBody}>
                  <EditPanel
                    doc={contentData}
                    type={type}
                    callbackContent={this.callbackContent}
                  />
                </Row>
              )}
              {isEdit === false && (
                <Row className={style.docBody}>
                  <div className={style.docContent}>
                    <ShowPanel doc={contentData} />
                  </div>
                </Row>
              )}
            </div>
          </div>
      </div>
    )
  }
}
