import React from 'react';
import { createStyles, Theme, makeStyles } from '@material-ui/core/styles'
import Image from '../TrueImages/background_Summary_16-9.png'
import { Box, CircularProgress, FormControlLabel, Switch } from '@material-ui/core'
import { Grid } from '@material-ui/core'
import { SessionData, getPatientSessions, getSessionsData } from '../funcs'
import GraphSessionComponent from '../GraphSessionComponent'
import { UserAppBar, UserInterfaceRole, UserType } from '../UserAppBar'

const useStyles = makeStyles((theme: Theme) => createStyles({
  background: {
    height: "100vh",
    width: "100vw",
    backgroundSize: 'cover'
  },
  buttonStyle: {
    background: "#cac7d6",
    marginRight: "8px"
  }
})
)

interface PatientSummaryProps { }

const PatientSummary = (props: PatientSummaryProps) => {
  const id: string = sessionStorage.getItem('id')!

  const classes = useStyles();

  const [AllSessionData, SetAllSessionData] = React.useState<SessionData[]>([])
  const [CompareSessions, SetCompareSessions] = React.useState<boolean>(false)
  const [CompareGraphs, SetCompareGraphs] = React.useState<boolean>(false)

  React.useEffect(() => {
    getPatientSessions(id).then(async (patientSessions) => {
      getSessionsData(patientSessions).then((value) => {
        SetAllSessionData(value)
      })
    })
  }, [])

  const changeCompareSessions = () => {
    SetCompareSessions(!CompareSessions)
  }

  const changeCompareGraphs = () => {
    SetCompareGraphs(!CompareGraphs)
  }

  return (
    <Box justifyContent="center"
      className={classes.background}
      style={{ backgroundImage: `url(${Image})` }}
    >
      <UserAppBar in={ UserInterfaceRole.Summary } for={ UserType.Patient } />

      <div style={{ padding: 16, marginTop: "32px" }}>
        {
          AllSessionData.length == 0 &&
          <CircularProgress />
        }
        {AllSessionData.length > 0 &&
          <Grid container>
            {
              !CompareSessions &&
              <Grid container>
                <Grid xs={7} />
                <Grid xs={2} >
                  <FormControlLabel
                    control={
                      <Switch checked={CompareGraphs} onChange={changeCompareGraphs} value="CompareGraphs" />
                    }

                    label={<span style={{fontSize:"42px"}}> Compare Graphs </span>}
                  />
                </Grid>
                <Grid xs={3} >
                  <FormControlLabel
                    control={
                      <Switch checked={CompareSessions} onChange={changeCompareSessions} value="CompareSessions" />
                    }
                    label={<span style={{fontSize:"42px"}}> Compare Sessions </span>}
                  />
                </Grid>
              </Grid>
            }
            {
              CompareSessions &&
              <Grid container>
                <Grid xs={10} />
                <Grid xs={2} >
                  <FormControlLabel
                    control={
                      <Switch checked={CompareSessions} onChange={changeCompareSessions} value="CompareSessions" />
                    }
                    label={<span style={{fontSize:"42px"}}> Compare Sessions </span>}
                  />
                </Grid>
              </Grid>
            }
            {
              !CompareSessions &&
              <Grid xs={12} >
                <GraphSessionComponent
                  SessionDataArr={AllSessionData}
                  CompareGraphs={CompareGraphs}
                />
              </Grid>
            }
            {
              CompareSessions &&
              <Grid container>
                <Grid xs={6} >
                  <GraphSessionComponent SessionDataArr={AllSessionData} />
                </Grid>
                <Grid xs={6} >
                  <GraphSessionComponent SessionDataArr={AllSessionData} />
                </Grid>
              </Grid>
            }
          </Grid>
        }
      </div>
    </Box>
  )
}

export default (PatientSummary)
