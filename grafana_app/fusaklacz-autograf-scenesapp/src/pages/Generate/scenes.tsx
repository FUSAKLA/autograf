import {
  DataSourceVariable,
  EmbeddedScene,
  TextBoxVariable,
  SceneControlsSpacer,
  SceneFlexLayout,
  SceneRefreshPicker,
  SceneTimePicker,
  SceneTimeRange,
  SceneVariableSet,
  VariableValueSelectors,
  SceneGridRow,
  PanelBuilders,
  SceneFlexItem,
  SceneGridLayout,
  SceneGridItem,
} from '@grafana/scenes';

export function getBasicScene(templatised = true, seriesToShow = '__server_names') {
  // Query runner definition, using Grafana built-in TestData datasource
  // const queryRunner = new SceneQueryRunner({
  //   datasource: datasourceVariable.getDataSourceRef(),
  //   queries: [
  //     {
  //       refId: 'A',
  //       datasource: DATASOURCE_REF,
  //       scenarioId: 'random_walk',
  //       seriesCount: 5,
  //       // Query is using variable value
  //       alias: templatised ? '${seriesToShow}' : seriesToShow,
  //       min: 30,
  //       max: 60,
  //     },
  //   ],
  //   maxDataPoints: 100,
  // });

  // // Query runner activation handler that will update query runner state when custom object state changes
  // queryRunner.addActivationHandler(() => {
  //   const sub = customObject.subscribeToState((newState) => {
  //     queryRunner.setState({
  //       queries: [
  //         {
  //           ...queryRunner.state.queries[0],
  //           seriesCount: newState.counter,
  //         },
  //       ],
  //     });
  //     queryRunner.runQueries();
  //   });

  //   return () => {
  //     sub.unsubscribe();
  //   };
  // });

  return new EmbeddedScene({
    $timeRange: new SceneTimeRange({
      from: 'now-1h',
      to: 'now',
    }),
    $variables: new SceneVariableSet({
      variables: [
        new DataSourceVariable({
          name: 'datasource',
          pluginId: 'prometheus',
        }),
        new TextBoxVariable({
          name: 'query',
        }),
      ],
    }),
    // $data: queryRunner,
    body: new SceneGridLayout({
      isLazy: true,
      children: [
        new SceneGridRow({
          x: 0,
          y: 0,
          title: 'Row 1',
          isCollapsed: false,
          isCollapsible: true,
          children: [
            new SceneGridItem({
              x: 0,
              y: 0,
              width: 12,
              height: 10,
              body: PanelBuilders.timeseries().setTitle('Time series').build(),
            }),
          ],
        }),
      ],
    }),
    controls: [
      new VariableValueSelectors({}),
      new SceneControlsSpacer(),
      new SceneTimePicker({ isOnCanvas: true }),
      new SceneRefreshPicker({
        intervals: ['5s', '1m', '1h'],
        isOnCanvas: true,
      }),
    ],
  });
}
