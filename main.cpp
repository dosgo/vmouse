#include <sys/types.h>
#include<pthread.h>
#include <stdio.h>
#include <string.h>
#include <windows.h>
#include <unistd.h>
#include <stdlib.h>
#include "cJSON.h"
#include "mouse.h"
int main(int argc, char **argv)
{
    if (argc != 2)
    {
        printf("Usage: %s port\n", argv[0]);
        exit(1);
    }
    printf("Welcome! virtual mouse！\n");


    #if WIN32
	WSADATA wsaData = { 0 };
	if ( 0 != WSAStartup( MAKEWORD( 2, 2 ), &wsaData ) )
	{
		printf( "WSAStartup failed. errno=[%d]\n", WSAGetLastError() );
		return(-1);
	}
    #endif


    struct sockaddr_in addr;
    addr.sin_family = AF_INET;
    addr.sin_port = htons(atoi(argv[1]));
    addr.sin_addr.s_addr = htonl(INADDR_ANY);

    int sock;
    if ( (sock = socket(AF_INET, SOCK_DGRAM, 0)) < 0)
    {
        perror("socket");
        exit(1);
    }
    if (bind(sock, (struct sockaddr *)&addr, sizeof(addr)) < 0)
    {
        perror("bind");
        exit(1);
    }
    char buff[512];
    struct sockaddr_in clientAddr;
    int n;
    int len = sizeof(clientAddr);
    while (1)
    {
        n = recvfrom(sock, buff, 511, 0, (struct sockaddr*)&clientAddr, &len);
        if (n>0)
        {
             cJSON *json = cJSON_Parse( buff );
             cJSON *cmd = cJSON_GetObjectItem( json, "cmd" );
             //ÍË³ö
             if ( strcmp( cmd->valuestring, "move" ) == 0 ){
                cJSON *x = cJSON_GetObjectItem( json, "x" );
                cJSON *y = cJSON_GetObjectItem( json, "y");
                if(x->type==cJSON_Number&&y->type==cJSON_Number){
                     mousemove(x->valueint,y->valueint);
                }

             }
             if ( strcmp( cmd->valuestring, "click" ) == 0 ){

                cJSON *type = cJSON_GetObjectItem( json, "type" );
                cJSON *doublec = cJSON_GetObjectItem( json, "double");
                if(type->type==cJSON_Number){
                     mouseclick(type->valueint,doublec->valueint);
                }

             }
            cJSON_Delete( json );
        }
        else
        {
            perror("recv");
            break;
        }
    }
    return 0;
}


