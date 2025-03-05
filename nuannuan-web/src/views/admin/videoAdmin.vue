<script setup lang="ts">
import { UploadFilled } from '@element-plus/icons-vue'
import router from "@/router";
import axios from "axios";
import {ElMessage, type InputInstance, type UploadFile, type UploadFiles, type UploadProgressEvent} from "element-plus";
import {nextTick, ref} from "vue";
import { ElInput } from 'element-plus'
const fileList = ref<UploadFile[]>([])

const goldFileList = ref<UploadFiles>([])

/**
 * tag
 */
const inputValue = ref('')
// const dynamicTags = ref(['Tag 1', 'Tag 2', 'Tag 3'])
const dynamicTags = ref<string[]>(['色色'])
const inputVisible = ref(false)
const InputRef = ref<InputInstance>()

// 添加文件事件 官方文档是想要我们用on-change	来实现添加 没有其它接口 只需要判断有没有加进去就可以了
const handleAdd = (file: UploadFile) => {
  // uid
  const existingFile = fileList.value.find((existing) => existing.uid === file.uid);

  if (existingFile) {
    // 如果找到了相同的文件，可以选择显示一个消息或者执行其他逻辑
    ElMessage.warning("文件已存在！");
    return // 不添加文件并返回 false
  }
  fileList.value.push(file)
  ElMessage.success("添加成功！");


}

const handleRemove = (file: UploadFile) => {
  const index = fileList.value.findIndex(existingFile => existingFile.uid === file.uid)



    fileList.value.splice(index, 1); // 使用splice方法删除指定索引处的元素
    ElMessage.success("删除成功！");

}

const handlePush = () => {
  const formData = new FormData();
  console.log(fileList.value)
  console.log("看看此时 fileList的长度 : ", fileList.value.length)
  for(let i = 0; i < fileList.value.length; i ++ ) {
    const file = fileList.value[i].raw ?? fileList.value[i].raw as unknown as File;
    const tmp_name = `${Date.now()}_${file.name}`
    formData.append('file', file, tmp_name);
    pushType(tmp_name)
    console.log("测试formData此时数据tmp_name : ", tmp_name)
  }

  axios.post('/api/video/push', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
      // 注意：通常不需要手动设置Content-Type，因为axios会根据FormData自动设置
    }
  }).then(response => {

    console.log('上传成功:', response.data);
    ElMessage.success("上传成功！")
  }).catch(error => {
    console.error('上传失败:', error);
    ElMessage.success("上传失败！")
  });
}

const pushType = (tmp_name: string) => {
  axios.post('/api/video/type/push',{
    types: dynamicTags.value,
    name: tmp_name
  }).then(response => {
    console.log('Types上传成功:', response.data);
    ElMessage.success("上传成功！")
  }).catch(error => {
    console.error('上传失败:', error);
  });
}

/**
 * tag
 * @param tag
 */

const handleClose = (tag: string) => {
  dynamicTags.value.splice(dynamicTags.value.indexOf(tag), 1)
}

const showInput = () => {
  inputVisible.value = true
  nextTick(() => {
    // 鼠标点击其他地方就不显示
    // InputRef.value!.input!.focus()
  })
}

const handleInputConfirm = () => {
  if (inputValue.value) {
    dynamicTags.value.push(inputValue.value)
  }
  inputVisible.value = false
  inputValue.value = ''
}

</script>

<template>



  <div>

    <el-upload
        class="upload-demo mt-3"
        drag
        action="https://run.mocky.io/v3/9d059bf9-4660-45f2-925d-ce80ad6c4d15"
        multiple
        :auto-upload="false"
        :on-change="handleAdd"
        :on-remove="handleRemove"
    >
      <el-icon class="el-icon--upload"><upload-filled /></el-icon>
      <div class="el-upload__text">
        Drop file here or <em>click to upload</em>
      </div>
      <template #tip>
        <div class="el-upload__tip">
          jpg/png files with a size less than 500kb
        </div>
      </template>
    </el-upload>
<!--  TAG-->



    <div class="flex gap-2">
      <el-tag
          class="mr-3"
          v-for="tag in dynamicTags"
          :key="tag"
          closable
          :disable-transitions="false"
          @close="handleClose(tag)"
      >
        {{ tag }}

      </el-tag>


      <div v-if="inputVisible" style="display: inline-block;">
        <el-input
            ref="InputRef"
            v-model="inputValue"
            style="width: 3.9rem"
            size="small"
            @keyup.enter="handleInputConfirm"
            @blur="handleInputConfirm"
        />

      </div>


      <el-button v-else class="button-new-tag" size="small" @click="showInput">
        + New Tag
      </el-button>
    </div>




    <el-button type="primary" @click="handlePush" class="mt-3">Upload</el-button>
  </div>

</template>


<style lang="scss" scoped>



</style>
