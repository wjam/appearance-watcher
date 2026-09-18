#import <Foundation/Foundation.h>
#include <stdint.h>

const char* GetCurrentAppearance() {
    NSString *style = [[NSUserDefaults standardUserDefaults] stringForKey:@"AppleInterfaceStyle"];
    return (style == nil) ? "Light" : [style UTF8String];
}
