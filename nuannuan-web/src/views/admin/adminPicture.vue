<script setup lang="ts">

import {Delete, Download, Plus, ZoomIn} from "@element-plus/icons-vue";
import {ref} from "vue";
import {ElMessage, type UploadFile} from "element-plus";
import axios from "axios";
import router from "@/router";




const dialogImageUrl = ref('')
const dialogVisible = ref(false)
const disabled = ref(false)
const smallVisible = ref(false)
// 初始化 fileList
const fileList = ref<UploadFile[]>([])
const sendType = ref('')

const input = ref<string[]>()



const handleRemove = (file: UploadFile) => {
  const index = fileList.value.findIndex(existingFile => existingFile.uid === file.uid)
  if(index === -1) {
    ElMessage("删除对象为空")
    return
  }


  if (index !== -1) {
    fileList.value.splice(index, 1); // 使用splice方法删除指定索引处的元素
    file.name = ""
    ElMessage.success("删除成功！再次点击删除缓存");
  }else {
    ElMessage.warning("未找到要删除的文件");
  }



}




/**
 * 改写
 * 计算图片大小
 * @param file
 */
const getImageDimensions = (file: File): Promise<{ width: number; height: number }> => {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.onload = () => {
      resolve({ width: img.width, height: img.height });
    };
    img.onerror = (error) => {
      reject(new Error(`Failed to load image: ${error}`));
    };

    // 创建一个URL，指向文件的内容
    const url = URL.createObjectURL(file);
    // console.log("看看 现在的图片是啥 ", url)
    // 设置图像的源为刚刚创建的URL
    img.src = url;
    // 在图像加载完成后，记得释放URL对象以避免内存泄漏
    // img.onload= () => {
    //   URL.revokeObjectURL(url);
    // };
  });
};

//
const handlePictureCardPreview = async (file: UploadFile) => {
  try {
    // 假设 UploadFile 类型兼容 File 类型，或者您可以从 UploadFile 中获取 File 对象
    const fileAsFile = file.raw ?? file as unknown as File;
    const dimensions = await getImageDimensions(fileAsFile);
    // console.log('Image dimensions:', dimensions);
    if(dimensions.height < 1024 && dimensions.width < 1024) {
      smallVisible.value = true;
    }
    else {
      smallVisible.value = false
    }
    dialogImageUrl.value = file.url!;
    dialogVisible.value = true;
  } catch (error) {
    console.error('Error getting image dimensions:', error);
    // 处理错误，例如显示错误消息
  }
};



const handleDownload = (file: UploadFile) => {
  console.log(file)
  const a = document.createElement('a');
  a.style.display = 'none';
  // 解决断言
  if(file.url) {
    a.href = file.url;
    a.download = file.name;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
  }
}

const addPicture = (file: UploadFile) => {
  if(fileList.value.find(existingFile => existingFile.uid === file.uid)) {
    console.log("添加失败，已经有了")
    return
  }
  console.log("添加中： ", file)
  fileList.value.push(file)

  // const url = URL.createObjectURL(file)
  // console.log("应该是失败了： ", url)
  /**
   * const fileAsFile = file.raw ?? file as unknown as File;
   */
  const fileAsFile = file.raw ?? file as unknown as File;
  const url = URL.createObjectURL(fileAsFile)
  console.log("现在生成的地址 ： ", url)


}

const handlePush = () => {

  const formData = new FormData();

  for(let i = 0; i < fileList.value.length; i ++ ) {

    /**
     * 生成图片
     */
        // 转换
    const file = fileList.value[i].raw ?? fileList.value[i].raw as unknown as File;
    const tmp_name = `${Date.now()}_${file.name}`
    /**
     * name：字段的名称，即表单中<input>元素的name属性值。
     * value：字段的值。这可以是字符串或Blob（包括子类型，如File）。
     * fileName（可选）：如果value是一个Blob对象，这个参数指定了文件的名称，即<input type="file">元素的filename属性值。如果未指定，则使用value的默认名称。
     */
    /**
     * 传多组的
     */
    formData.append('file', file, tmp_name);
    pushTag(tmp_name)

  }

  axios.post('/api/picture/push', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
      // 注意：通常不需要手动设置Content-Type，因为axios会根据FormData自动设置
    }
  }).then(response => {

    console.log('上传成功:', response.data);
    for(let i = 0; i < fileList.value.length; i ++ ) {
      handleRemove(fileList.value[i])
    }
    ElMessage.success("上传成功！")
  }).catch(error => {
    console.error('上传失败:', error);
  });

}

const pushTag = (tmp_name: string) => {
  console.log("看看 tag ： input value : ", input.value)
  axios.post('/api/tags/push',{

    tags: input.value,
    name: tmp_name
  }).then(response => {
    console.log('tags上传成功:', response.data);
    ElMessage.success("上传成功！")
  }).catch(error => {
    console.error('上传失败:', error);
  });
}





</script>

<template>



  <el-upload action="#" list-type="picture-card" :auto-upload="false">
    <el-icon><Plus /></el-icon>

    <template #file="{ file }">

      <div v-if="file.name">
        <img class="el-upload-list__item-thumbnail" :src="file.url" alt="" />
        <span class="el-upload-list__item-actions">
                      <span v-if="file.name">{{addPicture(file)}}</span>
          <span
              class="el-upload-list__item-preview"
              @click="handlePictureCardPreview(file)"
          >
            <el-icon><zoom-in /></el-icon>
          </span>
          <span
              v-if="!disabled"
              class="el-upload-list__item-delete"
              @click="handleDownload(file)"
          >
            <el-icon><Download /></el-icon>
          </span>
          <span
              v-if="!disabled"
              class="el-upload-list__item-delete"
              @click="handleRemove(file)"
          >
            <el-icon><Delete /></el-icon>
          </span>
        </span>
      </div>


    </template>


  </el-upload>


  <el-input-tag
      v-model="input"
      placeholder="Please input"
      aria-label="Please click the Enter key after input"
  />
  <el-button @click="handlePush">提交</el-button>


  <el-dialog v-model="dialogVisible" class="dialog-style">
    <img v-show="smallVisible" class="small-image-style" :src="dialogImageUrl" alt="Preview Image" />
    <img v-show="!smallVisible" class="big-image-style" :src="dialogImageUrl" alt="Preview Image" />
  </el-dialog>
</template>

<style scoped lang="scss">

</style>