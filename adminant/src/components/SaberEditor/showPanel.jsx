import React, {PureComponent} from 'react';
import {Form} from 'antd';


@Form.create()
export default class ShowPanel extends PureComponent {
  constructor(props) {
    super(props);
    this.showImg = false;
  }

  componentDidMount() {
    const that = this;
    const node = document.getElementById('editormd');
    // const $ = jquery;
    if (node) {
      $('#editormd').html('');
      // console.log("content", that.props.docBefore.content)
      editormd.markdownToHTML('editormd', {
        markdown: that.props.doc.preMarkdown, // + "\r\n" + $("#append-test").text(),
        htmlDecode: 'style,script,iframe', // you can filter tags decode
        toc             : true,
        tocm: true, // Using [TOCM]
        emoji: true,
        taskList: true,
        tex: true,
        flowChart: true,
        sequenceDiagram: true,
      });
    }
    $('img').css('cursor', 'pointer');
    $('img').bind('click', function () {
      const _this = $(this);
      if (that.showImg == false) {
        that.imgShow('#outerdiv', '#innerdiv', '#bigimg', _this, that);
        that.showImg = true;
      }
    });
  }
  imgShow(outerdiv, innerdiv, bigimg, _this, that) {
    const src = _this.attr('src');// 获取当前点击的pimg元素中的src属性
    $(bigimg).attr('src', src);// 设置#bigimg元素的src属性

    /* 获取当前点击图片的真实大小，并显示弹出层及大图 */
    $('<img/>').attr('src', src).load(function () {
      const windowW = $(window).width();// 获取当前窗口宽度
      const windowH = $(window).height();// 获取当前窗口高度
      const realWidth = this.width;// 获取图片真实宽度
      const realHeight = this.height;// 获取图片真实高度
      let imgWidth,
        imgHeight;
      const scale = 0.95;// 缩放尺寸，当图片真实宽度和高度大于窗口宽度和高度时进行缩放

      if (realHeight > windowH * scale) { // 判断图片高度
        imgHeight = windowH * scale;// 如大于窗口高度，图片高度进行缩放
        imgWidth = imgHeight / realHeight * realWidth;// 等比例缩放宽度
        if (imgWidth > windowW * scale) { // 如宽度扔大于窗口宽度
          imgWidth = windowW * scale;// 再对宽度进行缩放
        }
      } else if (realWidth > windowW * scale) { // 如图片高度合适，判断图片宽度
        imgWidth = windowW * scale;// 如大于窗口宽度，图片宽度进行缩放
        imgHeight = imgWidth / realWidth * realHeight;// 等比例缩放高度
      } else { // 如果图片真实高度和宽度都符合要求，高宽不变
        imgWidth = realWidth;
        imgHeight = realHeight;
      }
      $(bigimg).css('width', imgWidth);// 以最终的宽度对图片缩放

      const w = (windowW - imgWidth) / 2;// 计算图片与窗口左边距
      const h = (windowH - imgHeight) / 2;// 计算图片与窗口上边距
      $(innerdiv).css({ top: h, left: w });// 设置#innerdiv的top和left属性
      $(outerdiv).fadeIn('fast');// 淡入显示#outerdiv及.pimg
    });

    $(outerdiv).click(function () { // 再次点击淡出消失弹出层
      $(this).fadeOut('fast');
      that.showImg = false;
    });
  }
  componentWillReceiveProps(nextProps) {
    const idPre = this.props.doc.id || '';
    const idCur = nextProps.doc.id || '';
    if (idCur != '' && idCur != idPre) {
      const that = this;
      const content = nextProps.doc && nextProps.doc.preMarkdown || '';
      $('#editormd').html('');
      editormd.markdownToHTML('editormd', {
        markdown: content, // + "\r\n" + $("#append-test").text(),
        // htmlDecode      : true,       //
        htmlDecode: 'style,script,iframe', // you can filter tags decode
        toc: true,
        tocm: true, // Using [TOCM]
        // tocContainer    : "#custom-toc-container",
        // gfm             : false,
        // tocDropdown     : true,
        // markdownSourceCode : true,
        emoji: true,
        taskList: true,
        tex: true,
        flowChart: true,
        sequenceDiagram: true,
      });
      $('img').css('cursor', 'pointer');
      $('img').bind('click', function () {
        const _this = $(this);
        if (that.showImg == false) {
          that.imgShow('#outerdiv', '#innerdiv', '#bigimg', _this, that);
          that.showImg = true;
        }
      });
    }

  }
  render() {
    return (
      <div >
        <div id="editormd" />
        <div id="outerdiv" className="outerdiv">
          <div id="innerdiv" className="innerdiv">
            <img id="bigimg" className="bigimg" />
          </div>
        </div>
      </div>
    );
  }
}

