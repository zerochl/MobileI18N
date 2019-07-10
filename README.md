## MobileI18N
International document automatic generation tool for android and ios

android与ios端一键导出国际化文案工具

## 使用说明
1. key分三级（如有需要可自行调整代码）：page，element，element_key;android 按照“.”相连，ios使用“_”相连
2. cn.txt存放中文文案；en.txt存放英文文案
3.执行之后文案会输出在output目录

## 特殊符号说明

1. "&&"：百分号
2. "%@"：android：会被转义成%1$s,有多个会自动叠加序列，ios无影响
3. "\&\#160;"：空格
4. "\&\#8230;"：省略号
5. "&"：android：会被转移成"&#38;"，ios无影响

## 作者联系方式：QQ：975804495
## 疯狂的程序员群：186305789，没准你能遇到绝影大神
