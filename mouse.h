/*
  Copyright (c) 2009 Dave Gamble

  Permission is hereby granted, free of charge, to any person obtaining a copy
  of this software and associated documentation files (the "Software"), to deal
  in the Software without restriction, including without limitation the rights
  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  copies of the Software, and to permit persons to whom the Software is
  furnished to do so, subject to the following conditions:

  The above copyright notice and this permission notice shall be included in
  all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
  THE SOFTWARE.
*/

#ifndef mouse__h
#define mouse__h

#ifdef __cplusplus
extern "C"
{
#endif
//mousemove
inline int mousemove(int x, int y){
    HINSTANCE hDll;
    typedef bool (*Fun1)(int,int);
    hDll = LoadLibrary("user32.dll");
    if(NULL == hDll)
  {
    fprintf(stderr, "load dll 'user32.dll' fail.");
    return -1;
  }

  Fun1 SetCursorPos = (Fun1)GetProcAddress(hDll, "SetCursorPos");
  if(NULL == SetCursorPos)
  {
    fprintf(stderr, "call function 'SetCursorPos' fail.");
    FreeLibrary(hDll);
    return -1;
  }
  SetCursorPos(x,y);
  FreeLibrary(hDll);
  return 0;
}
//mouseclick

inline int mouseclick(int type,bool double_click){
    int left_click = MOUSEEVENTF_LEFTDOWN | MOUSEEVENTF_LEFTUP;
    int right_click = MOUSEEVENTF_RIGHTDOWN | MOUSEEVENTF_RIGHTUP;
    int clicktype;
    HINSTANCE hDll;
    typedef void (*Fun2)(
            DWORD dwFlags,        // motion and click options
            DWORD dx,             // horizontal position or change
            DWORD dy,             // vertical position or change
            DWORD dwData,         // wheel movement
            ULONG_PTR dwExtraInfo // application-defined information
    );

    hDll = LoadLibrary("user32.dll");
    if(NULL == hDll)
  {
    fprintf(stderr, "load dll 'user32.dll' fail.");
    return -1;
  }

  Fun2 mouse_event = (Fun2)GetProcAddress(hDll, "mouse_event");
  if(NULL == mouse_event)
  {
    fprintf(stderr, "call function 'mouse_event' fail.");
    FreeLibrary(hDll);
    return -1;
  }
  if(type==0)
    clicktype = left_click;
  else
    clicktype = right_click;
  mouse_event (clicktype, 0, 0, 0, 0 );
    FreeLibrary(hDll);
    if(double_click)
        mouseclick(type,false);
  return 0;
}




#ifdef __cplusplus
}
#endif


#endif
