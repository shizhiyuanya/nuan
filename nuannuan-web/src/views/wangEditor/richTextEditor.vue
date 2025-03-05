<template>
  <div style="border: 1px solid #ccc">
    <Toolbar
        style="border-bottom: 1px solid #ccc"
        :editor="editorRef"
        :defaultConfig="toolbarConfig"
        :mode="mode"
    />
    <Editor
        style="height: 500px; overflow-y: hidden;"
        v-model="valueHtml"
        :defaultConfig="editorConfig"
        :mode="mode"
        @onCreated="handleCreated"
    />
    <!-- 提交按钮 -->

  </div>
  <button @click="submitArticle">提交文章</button>
</template>

<script setup>
import '@wangeditor/editor/dist/css/style.css';
import {onBeforeUnmount, ref, shallowRef, onMounted} from 'vue';
import {Editor, Toolbar} from '@wangeditor/editor-for-vue';

const editorRef = shallowRef();
const valueHtml = ref('<p>hello</p>');

const toolbarConfig = {};
const editorConfig = {
  placeholder: '请输入内容...',
  MENU_CONF: {
    uploadImage: {
      server: 'http://localhost:8080/upload/image',
      fieldName: 'file',
      maxFileSize: 5 * 1024 * 1024, // 5MB
      async customUpload(file, insertFn) {
        const formData = new FormData();
        formData.append('file', file);

        try {
          const response = await fetch('http://localhost:9090/upload/image', {
            method: 'POST',
            body: formData,
          });
          const result = await response.json();
          if (result.url) {
            insertFn(result.url);
          }
        } catch (error) {
          console.error('图片上传失败:', error);
        }
      },
    },
  },
};

const mode = 'default';

onMounted(() => {
  setTimeout(() => {
    valueHtml.value = '<p>模拟 Ajax 异步设置内容</p>';
  }, 1500);
});

onBeforeUnmount(() => {
  const editor = editorRef.value;
  if (editor == null) return;
  editor.destroy();
});

const handleCreated = (editor) => {
  editorRef.value = editor;
};

// 提交文章
const submitArticle = async () => {
  try {
    const response = await fetch('http://localhost:9090/save', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({content: valueHtml.value}),
    });
    const result = await response.json();
    console.log('文章提交成功:', result);
  } catch (error) {
    console.error('文章提交失败:', error);
  }
};
</script>
