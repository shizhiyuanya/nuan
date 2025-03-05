import CryptoJS from 'crypto-js';

// 密钥（32 字节，对应 AES-256）
const aesKey = CryptoJS.enc.Utf8.parse('');

/**
 * 加密函数
 * @param plaintext 明文
 * @returns 加密后的 Base64 字符串
 */
export function encrypt(plaintext: string): string {
  // 生成随机初始化向量（IV）
  const iv = CryptoJS.lib.WordArray.random(16); // 16 字节 IV

  // 使用 AES-CFB 模式加密
  const encrypted = CryptoJS.AES.encrypt(plaintext, aesKey, {
    iv: iv,
    mode: CryptoJS.mode.CFB,
    padding: CryptoJS.pad.NoPadding, // Go 的 CFB 模式不需要填充
  });

  // 将 IV 和密文拼接，并转换为 Base64 字符串
  const combined = CryptoJS.enc.Base64.stringify(iv.concat(encrypted.ciphertext));
  return combined;
}

/**
 * 解密函数
 * @param ciphertext 加密后的 Base64 字符串
 * @returns 解密后的明文
 */
export function decrypt(ciphertext: string): string {
  // 解码 Base64 数据
  const combined = CryptoJS.enc.Base64.parse(ciphertext);

  // 提取 IV 和密文
  const iv = CryptoJS.lib.WordArray.create(combined.words.slice(0, 4)); // 前 16 字节是 IV
  const encrypted = CryptoJS.lib.WordArray.create(combined.words.slice(4)); // 剩余部分是密文

  // 使用 AES-CFB 模式解密
  const decrypted = CryptoJS.AES.decrypt( ciphertext, aesKey, {
    iv: iv,
    mode: CryptoJS.mode.CFB,
    padding: CryptoJS.pad.NoPadding, // Go 的 CFB 模式不需要填充
  });

  // 返回解密后的明文
  return decrypted.toString(CryptoJS.enc.Utf8);
}

// 示例
// const plaintext = 'Hello, this is a secret message!';
// const encrypted = encrypt(plaintext);
// console.log('Encrypted:', encrypted);

// const decrypted = decrypt(encrypted);
// console.log('Decrypted:', decrypted);