# vmouse
This is a virtual mouse set of services.


使用方法
启动服务
mouse.exe  8888  //8888为监听的udp端口


然后往本机监听的端口发送udp包就行了，包的格式



{"cmd":"move","x":500,"y":350} //移动鼠标到x坐标500 y坐标350出 mouse.ex收到了命令就会自动移动鼠标


{"cmd":"click","type":1,"double":0} //点击鼠标，type 0为鼠标左键 1为鼠标右键   double表示双击0为不双击，1为双击。就这么简单


use

Start the service

mouse.exe 8888 / / 8888  udp port


And then to the local monitoring port to send udp package on the line, the format of the package

{"cmd": "move", "x": 500, "y": 350} // move the mouse。

x=500  //Abscissa

y=350  //Y-axis
  

{"cmd":"click","type":1,"double":0}  // click the mouse

type=0  //left mouse  button

type=0  //right mouse  button

double=0 //double click
