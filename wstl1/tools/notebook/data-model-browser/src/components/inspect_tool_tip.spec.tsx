/**
 * @jsx React.createElement
 * @jsxFrag React.Fragment
 */
// Copyright 2020 Google LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// tslint:disable-next-line:enforce-name-casing
import * as React from 'react';
import {shallow} from 'enzyme';
import {InspectTooltip} from './inspect_tool_tip';

describe('InspectTooltip', () => {

  it('should render correctly', () => {
    const wrapper = shallow(
      <InspectTooltip title='Test Tool Tip'>
        <button id='testChild' type='submit'>
            Hello World
        </button>
      </InspectTooltip>);
    expect(wrapper).toMatchSnapshot();
  });
});
