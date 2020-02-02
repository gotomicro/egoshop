import React, { Component } from 'react';
import { connect } from 'dva';
import { message } from 'antd';

@connect(({ loading, editorModel }) => ({
  editorModel,
  loading: loading.models.editorModel,
}))
export default  class EditPanel extends Component {
    state = {
      doc: {},
      type: 0, //1 为feed，2 为com
      isDocChanged: true,
      visible: false,
    };
    constructor(props) {
      super(props);
      this.editor = {};
      this.hasLoaded = false;
    }

  componentDidMount() {
    const that = this;
    const node = document.getElementById('editormd');
    window.addEventListener('resize', () => { this.onWindowResize(); });
    if (node) {
      that.editor = editormd('editormd', {
        height: 600,
        path: `/static/editormd/lib/`,
        // path: '../src/static/editormd/lib/',
        placeholder: '本编辑器支持Markdown编辑，左边编写，右边预览',
        imageUpload: true,
        imageFormats: ['jpg', 'jpeg', 'gif', 'png', 'JPG', 'JPEG', 'GIF', 'PNG'],
        imageUploadURL: '/upload',
        taskList: true,
        tocm: true, // Using [TOCM]
        tex: true, // ¿ªÆô¿ÆÑ§¹«Ê½TeXÓïÑÔÖ§³Ö£¬Ä¬ÈÏ¹Ø±Õ
        flowChart: true, // ¿ªÆôÁ÷³ÌÍ¼Ö§³Ö£¬Ä¬ÈÏ¹Ø±Õ
        sequenceDiagram: true, // ¿ªÆôÊ±Ðò/ÐòÁÐÍ¼Ö§³Ö£¬Ä¬ÈÏ¹Ø±Õ,
        saveHTMLToTextarea: true,
        tocStartLevel: 1,
        toolbarIcons: ['save', 'undo', 'redo', 'h1', 'h2', 'h3', 'h4', 'bold', 'hr', 'italic', 'quote', 'list-ul', 'list-ol', 'link', 'reference-link', 'code', 'preformatted-text', 'code-block', 'table', 'watch', 'theme', '||', 'history', 'read',
          'full'],
        toolbarIconsClass: {
          bold: 'fa-bold',
        },
        toolbarIconTexts: {
          bold: 'a',
        },
        toolbarCustomIcons: {
          theme: '<a href="javascript:;" title="主题"> <i class="fa fa-tachometer" name="theme"></i></a>',
          save: '<a href="javascript:;" title="保存" id="markdown-save" class="disabled"> <i class="fa fa-save" name="save"></i></a>',
          history: '<a href="javascript:;" title="历史版本"> <i class="fa fa-history" name="history"></i></a>',
          read: '<a href="javascript:;" title="查看"> <i class="fa fa-eye" name="read"></i></a>',
          full: '<a href="javascript:;" title="全屏"> <i class="fa fa-arrows-alt" name="full"></i></a>',
        },
        toolbarHandlers: {
          /**
           * @param {Object}      cm         CodeMirror对象
           * @param {Object}      icon       图标按钮jQuery元素对象
           * @param {Object}      cursor     CodeMirror的光标对象，可获取光标所在行和位置
           * @param {String}      selection  编辑器选中的文本
           */
          save(cm, icon, cursor, selection) {
            if ($('#markdown-save').hasClass('change')) {
              that.saveDoc(true);
            }
          },
          history(cm, icon, cursor, selection) {
            that.showHistory();
          },
          theme(cm, icon, cursor, selection) {
            that.theme = !that.theme;
            if (that.theme) {
              that.editor.setTheme('dark');
              that.editor.setPreviewTheme('dark');
              that.editor.setEditorTheme('pastel-on-dark');
            } else {
              that.editor.setTheme('default');
              that.editor.setPreviewTheme('default');
              that.editor.setEditorTheme('default');
            }
          },
          read() {
            // that.props.history.push(`/resource/docshow?id=${}`)
            that.props.jumpRead();
          },
          home() {
            that.props.history.push(`/resource/docmanage`);
          },
        },
        onload() {
          const keyMap = {
            'Ctrl-S': function (cm) {
              if ($('#markdown-save').hasClass('change')) {
                that.saveDoc(true);
              }
            },
          };
          this.addKeyMap(keyMap);
          // 初始化文档
          if (that.props.doc && that.state.isDocChanged === true) {
            const content = that.props.doc && that.props.doc.preMarkdown || '';
            that.editor.clear();
            that.editor.insertValue(content);
            that.editor.setCursor({ line: 0, ch: 0 });
            that.setState({
              isDocChanged: true,
            });
            that.resetEditorChanged(false)
          }
          that.hasLoaded = true;
          that.timer1 = setInterval(() => {
            if ($('#markdown-save').hasClass('change')) {
              if (that.editor.getMarkdown().length > 0) {
                that.saveDoc(false);
              }
            }
          }, 10000);
          // that.timer2 = setInterval(() => {
          //   that.keepAlive();
          // }, 1000 * 60 * 10);
          that.editor.height(932);
          const winHeight = window.innerHeight;
          that.editor.height(winHeight - 10);
        },
        onchange() {
          that.resetEditorChanged(true)
          if (that.state.isDocChanged) {
            that.setState({
              isDocChanged: false,
            });
          } else {

          }
        },
      });
      // 图片上传
      $('#editormd').fileupload({
        dataType: 'json',
        pasteZone: $('#editormd'),
        dropZone: $('#editormd'), // 只允许paste
        paramName: 'file',
        progress(e, data) {
        },
        add(e, data) {
          const reader = new FileReader();
          reader.readAsDataURL(data.files[0]);
          reader.onload = function (e) {
            that.pasteImage({
              file: e.target.result,
            });
          };
        },
      });
    }
  }


  componentWillReceiveProps(nextProps) {
      this.setState({
        isDocChanged: true,
      });
      const that = this;

    const idPre = this.props.doc.id || '';
    const contentPre = this.props.doc.preMarkdown || '';
    const idCur = nextProps.doc.id || '';
    const contentCur = nextProps.doc.preMarkdown || '';
    // 当不id不相等的时候，触发
    if (idCur !== '' && idCur !== idPre) {
      if ((that.editor && that.hasLoaded) || (idCur == idPre && contentPre != contentCur)) {
          const content = nextProps.doc && nextProps.doc.preMarkdown || '';
          that.editor.clear();
          that.editor.insertValue(content);
          that.editor.setCursor({ line: 0, ch: 0 });
          $('#markdown-save').removeClass('change').addClass('disabled');

          that.editor.height(932);
          const winHeight = window.innerHeight;
          that.editor.height(winHeight - 10);
        }
    }
  }

  saveDoc(flag) {
    // console.log('hostname', window.location.hostname)
    const { dispatch,callbackContent } = this.props;
    const that = this;
    if (this.props.doc.id === '') {
      message.error('未选中要保存的文档');
      return;
    }
    dispatch({
      type: 'editorModel/contentSave',
      payload: {
        id: this.props.doc.id,
        type: this.props.type,
        markdown: that.editor.getMarkdown(),
        html: that.editor.getHTML(),
      },
      callback: (response) => {
        if (response.code === 0) {
            message.success('保存成功');
            callbackContent()
            $('#markdown-save').removeClass('change').addClass('disabled');
        } else {
          message.error(response.msg);
        }
      },
    });
  }

  pasteImage(data) {
    const that = this;
    const { dispatch } = this.props;
    dispatch({
      type: 'editorModel/upload',
      payload: {
        id:this.props.doc.id,
        type: this.props.type,
        image: data.file,
      },
      callback: (response) => {
        if (response.code === 0) {
          that.editor.insertValue(`![title](${response.data.url})`);
            $('#markdown-save').removeClass('change').addClass('disabled');
          } else {
            message.error(response.msg);

        }
      },
    });
  }

  /**
   * 设置编辑器变更状态
   * @param $is_change
   */
  resetEditorChanged($is_change) {
    if($is_change && !window.isLoad ){
      $("#markdown-save").removeClass('disabled').addClass('change');
    }else{
      $("#markdown-save").removeClass('change').addClass('disabled');
    }
    window.isLoad = false;
  }

  onWindowResize() {
    const winHeight = window.innerHeight;
    const that = this;
    if (that.editor && that.hasLoaded) {
        that.editor.height(winHeight - 10);
    }
  }
  render() {
    const that = this;
    let content = '';
    if (that.editor && that.hasLoaded) {
      content = that.editor.getMarkdown();
    }
    return (
      <div>
        <div name="editormd" id="editormd" />
      </div>
      )

  }
}
