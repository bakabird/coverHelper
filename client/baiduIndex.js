const net = require('net')
let times = 3

const startTime = (new Date).getTime()
var loop = setInterval(()=>{
    if(times-- >= 0){
        let step = 0
        const client = net.createConnection({
            port: 25005
        },()=>{
            console.log('connect suc')
        })
        urlsStr = `https://www.baidu.com/img/bd_logo1.png  http://i0.hdslb.com/bfs/archive/aa5ff915c3de108f4ca978e3e48b9aa3908ba40d.png https://www.baidu.com/img/bd_logo1.png https://www.baidu.com/img/bd_logo1.png http://i0.hdslb.com/bfs/archive/aa5ff915c3de108f4ca978e3e48b9aa3908ba40d.png http://i0.hdslb.com/bfs/archive/aa5ff915c3de108f4ca978e3e48b9aa3908ba40d.png`
        client.write(`C${urlsStr.length}`)
        client.on('data',(data)=>{
            console.log('recv data:',String(data) )
            // G = GOON
            if(String(data) === "G"){
                client.write(urlsStr)
            }else if(String(data).includes("OK")){
                client.write("R")
                step = 1
            }else if(step === 1){
                console.log("接收到转译结果为")
                console.log(String(data).split(' '))
                client.write("O")
                client.end()
            }else if(!String(data).includes("W")){
                console.log((new Date).getTime() - startTime)
                client.end()
            }
        })
    }else{
        clearInterval(loop)
    }
},100)
