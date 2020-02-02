import React, { Component, Fragment } from "react";
import { Select, Radio, InputNumber, Button } from "antd";
import router from "umi/router";
import { connect } from "dva";

const Option = Select.Option;
const RadioGroup = Radio.Group;

@connect(({ freightModel, loading }) => ({
    freightList: freightModel.list,
    freightListLoading: loading.effects["freightModel/list"]
}))
class GoodsFreight extends Component {
    static defaultProps = {
        freightList: [],
        freightListLoading: false,
        /**
         * 当freightId = 0时，运费算统一运费
         * @param e
         */
        onChange: (e = { freightType: "freightFee", freightFee: 0, freightId: 0 }) => {
        }
    };

    constructor(props) {
        super(props);
        const value = props.value || { freightId: 0, freightFee: 0 };
        this.state = {
            freightFee: value.freightFee,
            freightId: value.freightId,
            freightType: value.freightId > 0 ? "freightId" : "freightFee"
        };
        this.initList();
    }

    initList = () => {
        const { dispatch } = this.props;
        dispatch({
            type: "freight/list",
            payload: {
                page: 1,
                rows: 1000
            }
        });
    };

    componentWillReceiveProps(nextProps) {
        const value = nextProps.value || { freightId: 0, freightFee: 0 };
        this.setState({
            freightFee: value.freightFee,
            freightId: value.freightId,
            freightType: value.freightId > 0 ? "freightId" : "freightFee"
        })
    }

    render() {
        const { freightList, freightListLoading, onChange } = this.props;
        const { freightType, freightId, freightFee } = this.state;

        return (
            <RadioGroup
                onChange={(e) => {
                    this.setState({ freightType: e.target.value });
                }}
                value={freightType}
            >
                <Radio value={"freightFee"}>
                    统一邮费
                    {freightType === "freightFee" ? <InputNumber
                        placeholder="请输入"
                        style={{
                            width: 150,
                            marginLeft: 20
                        }}
                        formatter={value => `￥ ${value}`}
                        min={0}
                        precision={2}
                        value={freightFee}
                        onChange={(freightFee) => {
                            this.setState({ freightFee });
                            onChange({ freightType: "freightId", freightId: 0, freightFee });
                        }}
                    /> : null}
                </Radio>
                <br />
                <Radio value={"freightId"}>
                    运费模板
                    {freightType === "freightId" ? <Fragment><Select
                            showSearch
                            style={{ width: 240, marginLeft: 20 }}
                            placeholder="请选择"
                            optionFilterProp="children"
                            onChange={(value) => {
                                this.setState({
                                    freightId: value
                                });
                                onChange({ freightType: "freightId", freightId: value, freightFee: 0 });
                            }}
                            value={freightId === 0 ? null : freightId}
                            filterOption={(input, option) => option.props.children.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                        >
                            {
                                freightList.map((e) => (
                                    <Option value={e.id} key={e.id}>
                                        <p>{e.name}</p>
                                    </Option>
                                ))
                            }
                        </Select>
                            <Button
                                type="primary"
                                style={{
                                    marginLeft: 10
                                }}
                                size={"small"}
                                loading={freightListLoading}
                                onClick={this.initList}
                            >
                                刷新
                            </Button>
                        </Fragment>
                        : null}
                </Radio>
            </RadioGroup>
        );
    }
}

export default GoodsFreight;
