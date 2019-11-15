const net = require('net')
let times = 1

const startTime = (new Date).getTime()
var loop = setInterval(()=>{
    if(times-- >= 0){
        const client = net.createConnection({
            port: 25005
        },()=>{
            console.log('connect suc')
        })
        urlsStr = `https://www.baidu.com/img/bd_logo1.png`
        client.write(`C${urlsStr.length}`)
        client.on('data',(data)=>{
            console.log('recv data:',String(data) )
            // G = GOON
            if(String(data) === "G"){
                client.write(urlsStr)
            }else{
                console.log((new Date).getTime() - startTime)
                client.end()
            }
        })
    }else{
        clearInterval(loop)
    }
},1)
