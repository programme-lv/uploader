#include "testlib.h"
#include <iostream>

using namespace std;

int main(int argc, char* argv[])
{
    registerGen(argc, argv, 1);
    int n = atoi(argv[1]);
    cout << rnd.next(1, n) << " ";
    cout << rnd.next(1, n) << endl;
}
