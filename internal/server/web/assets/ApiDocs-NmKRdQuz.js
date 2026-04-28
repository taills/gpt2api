import{v as Oe,A as De,I as Ke,B as je,C as He,F as Se,i as qe,J as Fe,K as Ne,m as Ge,j as Je,f as Xe,L as Ze,u as Qe,x as We,M as et,N as tt,E as F,t as lt}from"./element-plus-D5RnCHG8.js";import{B as st,c as at,L as u,P as l,u as k,F as R,Y as n,V as N,Z as i,X as f,M as a,O as o,a as v,H as r,I as w,ab as B,a9 as ot,R as nt,f as y,e as it,ag as be}from"./vue-core-CfmJLmIH.js";import{l as rt,E as O,g as ut,a as dt,b as mt}from"./feature-CMatfP37.js";import{f as G,a as ke,b as ct}from"./format-CZ0sbVu3.js";import{_ as pt}from"./index-Ch8Mo1zj.js";import"./vendor-u8tqdTjN.js";const vt={class:"page-container"},_t={class:"card-block hero"},gt={class:"desc"},ft={class:"hero-stats"},yt={class:"stat"},bt={class:"val"},kt={key:0,class:"stat"},ht={class:"val"},wt={class:"stat"},Ct={class:"val"},$t={class:"stat"},Et={class:"val primary"},Pt={class:"card-block"},At={class:"row"},It={class:"code"},zt={class:"code"},xt={class:"card-block"},Rt={class:"flex-between",style:{"margin-bottom":"10px"}},Ut={key:0,class:"mute"},Tt={class:"pager"},Lt={class:"card-block"},Vt={class:"row"},Yt={class:"code"},Mt={class:"code"},Bt={class:"code"},Ot={class:"code"},Dt={class:"code"},Kt={class:"card-block"},jt={class:"flex-between",style:{"margin-bottom":"10px"}},Ht={key:0,class:"empty"},St={class:"grid"},qt=["onClick"],Ft=["src","alt"],Nt={key:1,class:"thumb-ph"},Gt={class:"s"},Jt={key:2,class:"thumb-badge"},Xt={class:"meta"},Zt=["title"],Qt={class:"sub"},Wt={class:"mute"},el={class:"foot"},tl={class:"mute"},ll={class:"credit"},sl={class:"actions"},al={key:0,class:"err"},ol={key:1,class:"pager"},nl={key:0},il=["title"],rl={class:"big-img-wrap"},ul={key:0,class:"thumb-strip"},dl=["src","onClick"],ml={class:"dlg-actions"},cl=st({__name:"ApiDocs",setup(pl){function J(s,e=10){if(!s)return s;const p=s.includes("?")?"&":"?";return`${s}${p}thumb_kb=${e}`}const X=v("image"),D=v([]),he=y(()=>D.value.filter(s=>s.type==="chat")),we=y(()=>D.value.filter(s=>s.type==="image")),U=v(""),h=v(""),C=y(()=>window.location.origin),E=v(null),K=v(!1);async function Ce(){K.value=!0;try{E.value=await ut({days:14,top_n:5})}finally{K.value=!1}}const Z=v([]),b=v({limit:20,offset:0,total:0}),j=v(!1);async function Q(){j.value=!0;try{const s=await mt({type:"chat",limit:b.value.limit,offset:b.value.offset});Z.value=s.items,b.value.total=s.total}finally{j.value=!1}}function $e(s){b.value.offset=(s-1)*b.value.limit,Q()}const T=v([]),P=v({limit:12,offset:0}),Y=v(!1),W=v(!1),c=it({status:"",keyword:"",range:[]});function Ee(){const s={};return c.status&&(s.status=c.status),c.keyword&&(s.keyword=c.keyword),c.range&&c.range.length===2&&(s.start_at=c.range[0],s.end_at=c.range[1]),s}async function A(s=!0){Y.value=!0;try{s&&(P.value.offset=0,T.value=[]);const e=await dt({limit:P.value.limit,offset:P.value.offset,...Ee()});s?T.value=e.items:T.value.push(...e.items),W.value=e.items.length>=P.value.limit}finally{Y.value=!1}}function Pe(){P.value.offset+=P.value.limit,A(!1)}function Ae(){c.status="",c.keyword="",c.range=[],A(!0)}const H=v(!1),I=v(null),z=v(0),M=y(()=>{var s;return((s=I.value)==null?void 0:s.image_urls)||[]}),ee=y(()=>M.value[z.value]||"");function te(s,e=0){var p;(p=s.image_urls)!=null&&p.length&&(I.value=s,z.value=e,H.value=!0)}async function le(s,e,p){if(s)try{const d=await fetch(s,{credentials:"include"});if(!d.ok)throw new Error("HTTP "+d.status);const m=await d.blob(),_=m.type||"image/png",L=_.includes("jpeg")?"jpg":_.split("/")[1]||"png",g=document.createElement("a"),V=URL.createObjectURL(m);g.href=V,g.download=`${e}-${p+1}.${L}`,document.body.appendChild(g),g.click(),document.body.removeChild(g),setTimeout(()=>URL.revokeObjectURL(V),6e4)}catch(d){F.error("下载失败:"+((d==null?void 0:d.message)||d))}}const se=y(()=>{const s=U.value||"gpt-5";return`curl ${C.value}/v1/chat/completions \\
  -H "Authorization: Bearer \${YOUR_API_KEY}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "${s}",
    "stream": true,
    "messages": [
      {"role": "user", "content": "你好,介绍一下你自己"}
    ]
  }'`}),ae=y(()=>{const s=U.value||"gpt-5";return`from openai import OpenAI

client = OpenAI(
    base_url="${C.value}/v1",
    api_key="\${YOUR_API_KEY}",
)

resp = client.chat.completions.create(
    model="${s}",
    messages=[{"role": "user", "content": "你好"}],
    stream=True,
)
for chunk in resp:
    print(chunk.choices[0].delta.content or "", end="")`}),oe=y(()=>{const s=h.value||"gpt-image-2";return`curl ${C.value}/v1/images/generations \\
  -H "Authorization: Bearer \${YOUR_API_KEY}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "${s}",
    "prompt": "A cute orange cat playing with yarn, studio ghibli style",
    "n": 1,
    "size": "1024x1024"
  }'`}),ne=y(()=>{const s=h.value||"gpt-image-2";return`curl ${C.value}/v1/images/generations \\
  -H "Authorization: Bearer \${YOUR_API_KEY}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "${s}",
    "prompt": "Restyle the cat as a watercolor painting, soft pastel palette",
    "n": 1,
    "size": "1024x1024",
    "reference_images": [
      "https://example.com/cat.png",
      "data:image/png;base64,iVBORw0KGgo..."
    ]
  }'`}),ie=y(()=>{const s=h.value||"gpt-image-2";return`import requests

API_KEY = "\${YOUR_API_KEY}"
BASE_URL = "${C.value}/v1"

resp = requests.post(
    f"{BASE_URL}/images/generations",
    headers={
        "Authorization": f"Bearer {API_KEY}",
        "Content-Type": "application/json",
    },
    json={
        "model": "${s}",
        "prompt": "A cute orange cat playing with yarn",
        "n": 1,
        "size": "1024x1024",
    },
    timeout=300,
)
resp.raise_for_status()
data = resp.json()
print(data["data"][0]["url"])`}),re=y(()=>{const s=h.value||"gpt-image-2";return`import base64, requests

API_KEY = "\${YOUR_API_KEY}"
BASE_URL = "${C.value}/v1"

# reference_images 单张最大 20MB,最多 4 张;支持 URL / data:URL / 纯 base64
def img_b64(path: str) -> str:
    with open(path, "rb") as f:
        return base64.b64encode(f.read()).decode()

resp = requests.post(
    f"{BASE_URL}/images/generations",
    headers={
        "Authorization": f"Bearer {API_KEY}",
        "Content-Type": "application/json",
    },
    json={
        "model": "${s}",
        "prompt": "Turn the cat into a watercolor painting",
        "n": 1,
        "size": "1024x1024",
        "reference_images": [
            img_b64("cat.png"),
            # 也可以直接传公网 URL:"https://example.com/style.jpg"
        ],
    },
    timeout=300,
)
resp.raise_for_status()
print(resp.json()["data"][0]["url"])`}),ue=y(()=>{const s=h.value||"gpt-image-2";return`from openai import OpenAI

client = OpenAI(
    base_url="${C.value}/v1",
    api_key="\${YOUR_API_KEY}",
)

resp = client.images.generate(
    model="${s}",
    prompt="A cute orange cat playing with yarn",
    n=1,
    size="1024x1024",
)
print(resp.data[0].url)`});async function $(s){try{await navigator.clipboard.writeText(s),F.success("已复制到剪贴板")}catch{F.error("复制失败,请手动选择文本")}}function de(s){return s==="success"?"success":s==="failed"?"danger":s==="running"||s==="dispatched"||s==="queued"?"warning":"info"}return at(async()=>{try{const s=await rt();D.value=O?s.items:s.items.filter(d=>d.type!=="chat");const e=s.items.find(d=>d.type==="chat"),p=s.items.find(d=>d.type==="image");e&&(U.value=e.slug),p&&(h.value=p.slug)}catch{}Ce(),A()}),(s,e)=>{var ce,pe,ve,_e,ge;const p=He,d=je,m=qe,_=Se,L=De,g=Ne,V=Ge,Ie=be("InfoFilled"),me=Xe,ze=Je,xe=Fe,Re=Ze,Ue=We,Te=et,Le=Qe,Ve=be("PictureRounded"),Ye=lt,Me=tt,Be=Ke,S=Oe;return r(),u("div",vt,[l("div",_t,[l("div",null,[e[27]||(e[27]=l("h2",{class:"page-title"},"接口文档 & 用量",-1)),l("p",gt,[k(O)?(r(),u(R,{key:0},[e[18]||(e[18]=n(" 外部调用走 ",-1)),e[19]||(e[19]=l("code",null,"/v1/chat/completions",-1)),e[20]||(e[20]=n(" 与 ",-1)),e[21]||(e[21]=l("code",null,"/v1/images/generations",-1)),e[22]||(e[22]=n(", ",-1))],64)):(r(),u(R,{key:1},[e[23]||(e[23]=n(" 外部调用走 ",-1)),e[24]||(e[24]=l("code",null,"/v1/images/generations",-1)),e[25]||(e[25]=n(", ",-1))],64)),e[26]||(e[26]=n(" 下面给出 curl / Python SDK 代码片段;个人用量与图片任务汇总在这里。若想在浏览器里直接体验,请打开「在线体验」。 ",-1))])]),N((r(),u("div",ft,[l("div",yt,[e[28]||(e[28]=l("div",{class:"lbl"},"14 天请求",-1)),l("div",bt,i(((ce=E.value)==null?void 0:ce.overall.requests)??0),1)]),k(O)?(r(),u("div",kt,[e[29]||(e[29]=l("div",{class:"lbl"},"文字 Token(in/out)",-1)),l("div",ht,i(((pe=E.value)==null?void 0:pe.overall.input_tokens)??0)+" / "+i(((ve=E.value)==null?void 0:ve.overall.output_tokens)??0),1)])):f("",!0),l("div",wt,[e[30]||(e[30]=l("div",{class:"lbl"},"图片张数",-1)),l("div",Ct,i(((_e=E.value)==null?void 0:_e.overall.image_images)??0),1)]),l("div",$t,[e[31]||(e[31]=l("div",{class:"lbl"},"14 天消耗积分",-1)),l("div",Et,i(k(G)((ge=E.value)==null?void 0:ge.overall.credit_cost)),1)])])),[[S,K.value]])]),a(L,{modelValue:X.value,"onUpdate:modelValue":e[15]||(e[15]=t=>X.value=t),class:"pg-tabs"},{default:o(()=>[k(O)?(r(),w(_,{key:0,label:"对话生成(文字模型)",name:"chat"},{default:o(()=>[l("div",Pt,[l("div",At,[e[32]||(e[32]=l("div",{class:"label"},"文字模型",-1)),a(d,{modelValue:U.value,"onUpdate:modelValue":e[0]||(e[0]=t=>U.value=t),placeholder:"选择模型",style:{width:"320px"}},{default:o(()=>[(r(!0),u(R,null,B(he.value,t=>(r(),w(p,{key:t.id,label:`${t.slug}${t.description?" · "+t.description:""}`,value:t.slug},null,8,["label","value"]))),128))]),_:1},8,["modelValue"])]),a(L,{type:"border-card",class:"code-tabs"},{default:o(()=>[a(_,{label:"curl"},{default:o(()=>[l("pre",It,[l("code",null,i(se.value),1)]),a(m,{size:"small",onClick:e[1]||(e[1]=t=>$(se.value))},{default:o(()=>[...e[33]||(e[33]=[n("复制 curl",-1)])]),_:1})]),_:1}),a(_,{label:"Python (OpenAI SDK)"},{default:o(()=>[l("pre",zt,[l("code",null,i(ae.value),1)]),a(m,{size:"small",onClick:e[2]||(e[2]=t=>$(ae.value))},{default:o(()=>[...e[34]||(e[34]=[n("复制 Python",-1)])]),_:1})]),_:1})]),_:1})]),l("div",xt,[l("div",Rt,[e[36]||(e[36]=l("h3",{class:"section-title"},"文字调用历史",-1)),a(m,{size:"small",onClick:Q},{default:o(()=>[...e[35]||(e[35]=[n("刷新",-1)])]),_:1})]),N((r(),w(xe,{data:Z.value,stripe:"",size:"small"},{default:o(()=>[a(g,{prop:"created_at",label:"时间","min-width":"160"},{default:o(({row:t})=>[n(i(k(ke)(t.created_at)),1)]),_:1}),a(g,{prop:"model_slug",label:"模型","min-width":"140"}),a(g,{label:"Token (in / out / cache)","min-width":"170"},{default:o(({row:t})=>[n(i(t.input_tokens)+" / "+i(t.output_tokens)+" ",1),t.cache_read_tokens?(r(),u("span",Ut,"/ "+i(t.cache_read_tokens),1)):f("",!0)]),_:1}),a(g,{label:"耗时",width:"90"},{default:o(({row:t})=>[n(i(t.duration_ms)+" ms",1)]),_:1}),a(g,{label:"状态",width:"90"},{default:o(({row:t})=>[a(V,{type:de(t.status),size:"small"},{default:o(()=>[n(i(t.status),1)]),_:2},1032,["type"]),t.error_code?(r(),w(ze,{key:0,content:k(ct)(t.error_code)+"("+t.error_code+")"},{default:o(()=>[a(me,{style:{"margin-left":"4px"}},{default:o(()=>[a(Ie)]),_:1})]),_:1},8,["content"])):f("",!0)]),_:1}),a(g,{label:"扣费(积分)",width:"110"},{default:o(({row:t})=>[n(i(k(G)(t.credit_cost)),1)]),_:1})]),_:1},8,["data"])),[[S,j.value]]),l("div",Tt,[a(Re,{layout:"prev, pager, next, total",total:b.value.total,"page-size":b.value.limit,"current-page":Math.floor(b.value.offset/b.value.limit)+1,onCurrentChange:$e},null,8,["total","page-size","current-page"])])])]),_:1})):f("",!0),a(_,{label:"图片生成(图片模型)",name:"image"},{default:o(()=>[l("div",Lt,[l("div",Vt,[e[37]||(e[37]=l("div",{class:"label"},"图片模型",-1)),a(d,{modelValue:h.value,"onUpdate:modelValue":e[3]||(e[3]=t=>h.value=t),placeholder:"选择模型",style:{width:"320px"}},{default:o(()=>[(r(!0),u(R,null,B(we.value,t=>(r(),w(p,{key:t.id,label:`${t.slug}${t.description?" · "+t.description:""}`,value:t.slug},null,8,["label","value"]))),128))]),_:1},8,["modelValue"])]),a(L,{type:"border-card",class:"code-tabs"},{default:o(()=>[a(_,{label:"curl(文生图)"},{default:o(()=>[l("pre",Yt,[l("code",null,i(oe.value),1)]),a(m,{size:"small",onClick:e[4]||(e[4]=t=>$(oe.value))},{default:o(()=>[...e[38]||(e[38]=[n("复制 curl",-1)])]),_:1})]),_:1}),a(_,{label:"curl(图生图)"},{default:o(()=>[l("pre",Mt,[l("code",null,i(ne.value),1)]),e[40]||(e[40]=l("div",{class:"hint"},[n(" reference_images 支持 "),l("code",null,"URL"),n(" / "),l("code",null,"data:URL"),n(" / 纯 "),l("code",null,"base64"),n("; 单次最多 4 张,单张最大 20MB。 ")],-1)),a(m,{size:"small",onClick:e[5]||(e[5]=t=>$(ne.value))},{default:o(()=>[...e[39]||(e[39]=[n("复制 curl",-1)])]),_:1})]),_:1}),a(_,{label:"Python (OpenAI SDK)"},{default:o(()=>[l("pre",Bt,[l("code",null,i(ue.value),1)]),a(m,{size:"small",onClick:e[6]||(e[6]=t=>$(ue.value))},{default:o(()=>[...e[41]||(e[41]=[n("复制 Python",-1)])]),_:1})]),_:1}),a(_,{label:"Python (requests · 文生图)"},{default:o(()=>[l("pre",Ot,[l("code",null,i(ie.value),1)]),a(m,{size:"small",onClick:e[7]||(e[7]=t=>$(ie.value))},{default:o(()=>[...e[42]||(e[42]=[n("复制 Python",-1)])]),_:1})]),_:1}),a(_,{label:"Python (requests · 图生图)"},{default:o(()=>[l("pre",Dt,[l("code",null,i(re.value),1)]),e[44]||(e[44]=l("div",{class:"hint"},[n(" reference_images 同时支持 "),l("code",null,"URL"),n(" / "),l("code",null,"data:URL"),n(" / 纯 "),l("code",null,"base64"),n(", 最多 4 张、单张最大 20MB,服务端会自动下载并解码。 ")],-1)),a(m,{size:"small",onClick:e[8]||(e[8]=t=>$(re.value))},{default:o(()=>[...e[43]||(e[43]=[n("复制 Python",-1)])]),_:1})]),_:1})]),_:1})]),l("div",Kt,[l("div",jt,[e[46]||(e[46]=l("h3",{class:"section-title"},"图片任务历史",-1)),a(m,{size:"small",onClick:e[9]||(e[9]=t=>A(!0))},{default:o(()=>[...e[45]||(e[45]=[n("刷新",-1)])]),_:1})]),a(Le,{inline:"",class:"flex-wrap-gap",style:{"margin-bottom":"10px"},onSubmit:e[14]||(e[14]=ot(t=>A(!0),["prevent"]))},{default:o(()=>[a(Ue,{modelValue:c.keyword,"onUpdate:modelValue":e[10]||(e[10]=t=>c.keyword=t),placeholder:"提示词关键字",clearable:"",style:{width:"220px"}},null,8,["modelValue"]),a(d,{modelValue:c.status,"onUpdate:modelValue":e[11]||(e[11]=t=>c.status=t),placeholder:"状态",clearable:"",style:{width:"130px"}},{default:o(()=>[a(p,{label:"成功",value:"success"}),a(p,{label:"失败",value:"failed"}),a(p,{label:"运行中",value:"running"}),a(p,{label:"队列中",value:"queued"})]),_:1},8,["modelValue"]),a(Te,{modelValue:c.range,"onUpdate:modelValue":e[12]||(e[12]=t=>c.range=t),type:"datetimerange","unlink-panels":"","range-separator":"~","start-placeholder":"开始时间","end-placeholder":"结束时间",format:"YYYY-MM-DD HH:mm","value-format":"YYYY-MM-DD HH:mm:ss",style:{width:"340px"}},null,8,["modelValue"]),a(m,{type:"primary",onClick:e[13]||(e[13]=t=>A(!0))},{default:o(()=>[...e[47]||(e[47]=[n("查询",-1)])]),_:1}),a(m,{onClick:Ae},{default:o(()=>[...e[48]||(e[48]=[n("重置",-1)])]),_:1})]),_:1}),N((r(),u("div",null,[T.value.length===0&&!Y.value?(r(),u("div",Ht," 暂无图片任务,复制上方代码调用一次即可生成记录。 ")):f("",!0),l("div",St,[(r(!0),u(R,null,B(T.value,t=>(r(),w(Ye,{key:t.id,shadow:"hover",class:"img-card"},{default:o(()=>{var x,q,fe;return[l("div",{class:"thumb",onClick:ye=>te(t,0)},[(x=t.image_urls)!=null&&x[0]?(r(),u("img",{key:0,src:J(t.image_urls[0]),alt:t.prompt,loading:"lazy"},null,8,Ft)):(r(),u("div",Nt,[a(me,{size:32},{default:o(()=>[a(Ve)]),_:1}),l("div",Gt,i(t.status),1)])),t.image_urls&&t.image_urls.length>1?(r(),u("div",Jt,i(t.image_urls.length)+" 张 ",1)):f("",!0)],8,qt),l("div",Xt,[l("div",{class:"title",title:t.prompt},i(t.prompt||"(无 prompt)"),9,Zt),l("div",Qt,[a(V,{size:"small",type:de(t.status)},{default:o(()=>[n(i(t.status),1)]),_:2},1032,["type"]),l("span",null,i(t.size),1),l("span",Wt,"n="+i(t.n),1)]),l("div",el,[l("span",tl,i(k(ke)(t.created_at)),1),l("span",ll,i(k(G)(t.credit_cost))+" 积分",1)]),l("div",sl,[(q=t.image_urls)!=null&&q.length?(r(),w(m,{key:0,size:"small",type:"primary",link:"",onClick:ye=>te(t,0)},{default:o(()=>[...e[49]||(e[49]=[n("放大",-1)])]),_:1},8,["onClick"])):f("",!0),(fe=t.image_urls)!=null&&fe.length?(r(),w(m,{key:1,size:"small",link:"",onClick:ye=>le(t.image_urls[0],t.task_id,0)},{default:o(()=>[...e[50]||(e[50]=[n("下载",-1)])]),_:1},8,["onClick"])):f("",!0)]),t.error?(r(),u("div",al,i(t.error),1)):f("",!0)])]}),_:2},1024))),128))]),W.value?(r(),u("div",ol,[a(m,{onClick:Pe},{default:o(()=>[...e[51]||(e[51]=[n("加载更多",-1)])]),_:1})])):f("",!0)])),[[S,Y.value]])])]),_:1})]),_:1},8,["modelValue"]),a(Be,{modelValue:H.value,"onUpdate:modelValue":e[17]||(e[17]=t=>H.value=t),title:"图片预览",width:"780px"},{default:o(()=>[I.value?(r(),u("div",nl,[l("div",{class:"prompt-line",title:I.value.prompt},i(I.value.prompt),9,il),l("div",rl,[a(Me,{src:ee.value,"preview-src-list":M.value,"initial-index":z.value,fit:"contain",style:{"max-height":"60vh","max-width":"100%",cursor:"zoom-in"}},null,8,["src","preview-src-list","initial-index"])]),M.value.length>1?(r(),u("div",ul,[(r(!0),u(R,null,B(M.value,(t,x)=>(r(),u("img",{key:x,src:J(t,16),alt:"",loading:"lazy",class:nt(["p-thumb",{active:x===z.value}]),onClick:q=>z.value=x},null,10,dl))),128))])):f("",!0),l("div",ml,[a(m,{size:"small",type:"primary",onClick:e[16]||(e[16]=t=>le(ee.value,I.value.task_id,z.value))},{default:o(()=>[...e[52]||(e[52]=[n("下载当前",-1)])]),_:1})])])):f("",!0)]),_:1},8,["modelValue"])])}}}),kl=pt(cl,[["__scopeId","data-v-ef4cea96"]]);export{kl as default};
