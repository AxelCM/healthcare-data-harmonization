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
import {SubEntitiesViewer} from './sub_entities_viewer';
import {SubEntityMeta} from './types';

describe('SubEntitiesViewer', () => {
  it('should render correctly', () => {
    const entities: SubEntityMeta[] = [
      {
        name: 'Test1',
        required: true,
        type: 'string',
        description: 'Verbose details',
        clickable: true,
      }
    ];

    const selected = 'Test1';
    const onInspect = jest.fn();
    const onSelect = jest.fn();
    const wrapper = shallow(
      <SubEntitiesViewer
        subEntities={entities}
        onInspect={onInspect}
        onSelect={onSelect}
        selected={selected}
      />
    );
    expect(wrapper).toMatchSnapshot();
  });
});
