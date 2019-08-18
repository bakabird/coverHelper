const net = require('net')
let times = 100

const startTime = (new Date).getTime()
var loop = setInterval(()=>{
    if(times-- >= 0){
        const client = net.createConnection({
            port: 25003
        },()=>{
            console.log('connect suc')
        })
        urlsStr = `http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg http://i2.hdslb.com/bfs/archive/6ba5f7e15b9322e2b3653538bfcc614f657816e8.jpg`
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
