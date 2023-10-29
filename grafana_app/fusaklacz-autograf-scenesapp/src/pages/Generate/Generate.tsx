import React, { useMemo } from 'react';

import { SceneApp, SceneAppPage } from '@grafana/scenes';
import { getBasicScene } from './scenes';
import { prefixRoute } from '../../utils/utils.routing';
import { ROUTES } from '../../constants';
import { config } from '@grafana/runtime';
import { Alert } from '@grafana/ui';

const getScene = () => {
  return new SceneApp({
    pages: [
      new SceneAppPage({
        title: 'Autograf',
        subTitle: 'Automatically generate dashboard from Prometheus metrics',
        url: prefixRoute(ROUTES.Generate),
        getScene: () => {
          return getBasicScene();
        },
      }),
    ],
  });
};
export const HomePage = () => {
  const scene = useMemo(() => getScene(), []);

  return (
    <>
      {!config.featureToggles.topnav && (
        <Alert title="Missing topnav feature toggle">
          Scenes are designed to work with the new navigation wrapper that will be standard in Grafana 10
        </Alert>
      )}
      <scene.Component model={scene} />
    </>
  );
};
