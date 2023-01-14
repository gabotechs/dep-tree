// @ts-ignore
import { a, b } from './2/2'
import * as something from './2/index'
import './1/a/a'
// @ts-ignore
import('./1/a')
// @ts-ignore
import('./unexisting')
// @ts-ignore
import { Unexisting } from './1/a'
