export interface TotpEntry {
  id: number;
  account: string;
  issuer: string;
  secret: string;
  otpauth: string;
  remark: string;
  currentCode: string;
  remainSeconds: number;
  category: string;  // 新增字段
}
