import * as React from 'react';
import { Redirect, Route, Switch } from 'react-router-dom';
import { HomePage } from '../../pages/Generate';
import { prefixRoute } from '../../utils/utils.routing';
import { ROUTES } from '../../constants';

export const Routes = () => {
  return (
    <Switch>
      <Route path={prefixRoute(`${ROUTES.Generate}`)} component={HomePage} />
      <Redirect to={prefixRoute(ROUTES.Generate)} />
    </Switch>
  );
};
